import * as fs from 'fs';
import timeSpan from 'time-span';
import {
  Blocks,
  stakingWithoutSlashing,
  bondBasedConsumerVotingPower,
} from './properties.js';
import {
  NUM_VALIDATORS,
  P,
  C,
  MAX_VALIDATORS,
  TOKEN_SCALAR,
  BLOCK_SECONDS,
  SLASH_DOUBLESIGN,
  SLASH_DOWNTIME,
} from './constants.js';
import _ from 'underscore';
import { Model } from './model.js';
import { Event } from './events.js';
import { strict as assert } from 'node:assert';

function forceMakeEmptyDir(dir) {
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir);
    return;
  }
  fs.rmSync(dir, { recursive: true });
  forceMakeEmptyDir(dir);
}

interface Action {
  kind: string;
}

type Delegate = {
  kind: string;
  val: number;
  amt: number;
};
type Undelegate = {
  kind: string;
  val: number;
  amt: number;
};
type JumpNBlocks = {
  kind: string;
  chains: string[];
  n: number;
  secondsPerBlock: number;
};
type Deliver = {
  kind: string;
  chain: string;
};
type ProviderSlash = {
  kind: string;
  val: number;
  power: number;
  infractionHeight: number;
  factor: number;
};
type ConsumerSlash = {
  kind: string;
  val: number;
  power: number;
  infractionHeight: number;
  isDowntime: number;
};

function weightedRandomKey(distr) {
  const scalar = _.reduce(_.values(distr), (sum, y) => sum + y, 0);
  const x = Math.random() * scalar;
  const pairs = _.pairs(distr);
  let i = 0;
  let cum = 0;
  while (i < pairs.length - 1 && cum + pairs[i][1] < x) {
    cum += pairs[i][1];
    i += 1;
  }
  return pairs[i][0];
}

class ActionGenerator {
  model;
  delegatedSinceBlock = new Array(NUM_VALIDATORS).fill(false);
  undelegatedSinceBlock = new Array(NUM_VALIDATORS).fill(false);
  jailed = new Array(NUM_VALIDATORS).fill(false);
  lastJumped = [];

  constructor(model) {
    this.model = model;
  }

  get = () => {
    let templates: Action[] = _.flatten([
      this.candidateDelegate(),
      this.candidateUndelegate(),
      this.candidateJumpNBlocks(),
      this.candidateDeliver(),
      this.candidateProviderSlash(),
      this.candidateConsumerSlash(),
    ]);
    const possible = _.uniq(templates.map((a) => a.kind));
    const distr = _.pick(
      {
        Delegate: 0.03,
        Undelegate: 0.03,
        JumpNBlocks: 0.45,
        Deliver: 0.45,
        ProviderSlash: 0.02,
        ConsumerSlash: 0.02,
      },
      ...possible,
    );
    const kind = weightedRandomKey(distr);
    templates = templates.filter((a) => a.kind === kind);
    const a = _.sample(templates);
    if (kind === 'Delegate') {
      return this.selectDelegate(a);
    }
    if (kind === 'Undelegate') {
      return this.selectUndelegate(a);
    }
    if (kind === 'JumpNBlocks') {
      return this.selectJumpNBlocks(a);
    }
    if (kind === 'Deliver') {
      return this.selectDeliver(a);
    }
    if (kind === 'ProviderSlash') {
      return this.selectProviderSlash(a);
    }
    if (kind === 'ConsumerSlash') {
      return this.selectConsumerSlash(a);
    }
    throw 'invalid kind';
  };

  candidateDelegate = (): Action[] => {
    return _.range(NUM_VALIDATORS)
      .filter((i) => !this.delegatedSinceBlock[i])
      .map((i) => {
        return {
          kind: 'Delegate',
          val: i,
        };
      });
  };

  candidateUndelegate = (): Action[] => {
    return _.range(NUM_VALIDATORS)
      .filter((i) => !this.undelegatedSinceBlock[i])
      .map((i) => {
        return {
          kind: 'Undelegate',
          val: i,
        };
      });
  };

  candidateJumpNBlocks = (): Action[] => [{ kind: 'JumpNBlocks' }];

  candidateDeliver = (): Action[] => {
    return [P, C]
      .filter((c) => this.model.hasUndelivered(c))
      .map((c) => {
        return { kind: 'Deliver', chain: c };
      });
  };

  candidateProviderSlash = (): Action[] => {
    return _.range(NUM_VALIDATORS)
      .filter((i) => {
        const cntWouldBeNotJailed = this.jailed.filter(
          (j) => j !== i && !this.jailed[j],
        ).length;
        return MAX_VALIDATORS <= cntWouldBeNotJailed;
      })
      .map((i) => {
        return { kind: 'ProviderSlash', val: i };
      });
  };

  candidateConsumerSlash = (): Action[] => {
    return _.range(NUM_VALIDATORS)
      .filter((i) => {
        const cntWouldBeNotJailed = this.jailed.filter(
          (j) => j !== i && !this.jailed[j],
        ).length;
        return MAX_VALIDATORS <= cntWouldBeNotJailed;
      })
      .map((i) => {
        return { kind: 'ConsumerSlash', val: i };
      });
  };

  selectDelegate = (a): Delegate => {
    this.delegatedSinceBlock[a.val] = true;
    return { ...a, amt: _.random(1, 5) * TOKEN_SCALAR };
  };

  selectUndelegate = (a): Undelegate => {
    this.undelegatedSinceBlock[a.val] = true;
    return { ...a, amt: _.random(1, 4) * TOKEN_SCALAR };
  };

  selectJumpNBlocks = (a): JumpNBlocks => {
    const chains = _.sample([[P], [C], [P, C]]); //TODO:
    a = {
      ...a,
      chains,
      n: _.sample([1, 6]),
      secondsPerBlock: BLOCK_SECONDS,
    };
    if (a.chains.includes(P)) {
      this.delegatedSinceBlock = new Array(NUM_VALIDATORS).fill(false);
      this.undelegatedSinceBlock = new Array(NUM_VALIDATORS).fill(false);
    }
    return a;
  };
  selectDeliver = (a): Deliver => {
    return a;
  };
  selectProviderSlash = (a): ProviderSlash => {
    // TODO: can only happen with evidence, power can't be random
    this.jailed[a.val] = true;
    return {
      ...a,
      power: _.random(1, 6) * TOKEN_SCALAR,
      infractionHeight: Math.floor(Math.random() * this.model.h[P]),
      factor: _.sample([SLASH_DOUBLESIGN, SLASH_DOWNTIME]),
    };
  };
  selectConsumerSlash = (a): ConsumerSlash => {
    // TODO: can only happen with evidence, power can't be random
    this.jailed[a.val] = true;
    return {
      ...a,
      power: _.random(1, 6) * TOKEN_SCALAR,
      infractionHeight: Math.floor(Math.random() * this.model.h[C]),
      isDowntime: _.sample([true, false]),
    };
  };
}

class Trace {
  actions = [];
  consequences = [];
  blocks = [];
  events = [];
  dump = (fn: string) => {
    const transitions = _.zip(this.actions, this.consequences).map(
      ([a, c]) => {
        return { action: a, consequence: c };
      },
    );
    const json = JSON.stringify({ transitions }, null, 4);
    this.write(fn, json);
  };
  write = (fn, content) => {
    // fs.writeFile(fn, content, 'utf8', (err) => {
    // if (err) throw err;
    // });
    fs.writeFileSync(fn, content);
  };
}

function doAction(model, action: Action) {
  const kind = action.kind;
  if (kind === 'Delegate') {
    const a = action as Delegate;
    model.delegate(a.val, a.amt);
  }
  if (kind === 'Undelegate') {
    const a = action as Undelegate;
    model.undelegate(a.val, a.amt);
  }
  if (kind === 'JumpNBlocks') {
    const a = action as JumpNBlocks;
    model.jumpNBlocks(a.n, a.chains, a.secondsPerBlock);
  }
  if (kind === 'Deliver') {
    const a = action as Deliver;
    model.deliver(a.chain);
  }
  if (kind === 'ProviderSlash') {
    const a = action as ProviderSlash;
    model.providerSlash(a.val, a.infractionHeight, a.power, a.factor);
  }
  if (kind === 'ConsumerSlash') {
    const a = action as ConsumerSlash;
    model.consumerSlash(a.val, a.power, a.infractionHeight, a.isDowntime);
  }
  throw 'wrong kind';
}

function writeEventData(allEvents, fn) {
  const eventCnt = _.countBy(allEvents, _.identity);
  for (const evt in Event) {
    const evtName = Event[evt];
    if (!_.has(eventCnt, evtName)) {
      eventCnt[evtName] = 0;
    }
  }
  const cnts = _.chain(eventCnt)
    .pairs()
    .sortBy((pair) => -pair[1]);

  fs.writeFileSync(`cnts${fn}.json`, JSON.stringify(cnts));
}

function gen() {
  const outerEnd = timeSpan();
  const GOAL_TIME_MINS = 5;
  const goalTimeMillis = GOAL_TIME_MINS * 60 * 1000;
  const NUM_ACTIONS = 40;
  const DIR = 'traces/';
  forceMakeEmptyDir(DIR);
  let numRuns = 1000000000000;
  let elapsedMillis = 0;
  let i = 0;
  const allEvents = [];
  while (i < numRuns) {
    i += 1;
    numRuns = Math.round(goalTimeMillis / (elapsedMillis / i) + 0.5);
    const end = timeSpan();
    ////////////////////////
    const blocks = new Blocks();
    const events = [];
    const model = new Model(blocks, events);
    const actionGenerator = new ActionGenerator(model);
    const trace = new Trace();
    for (let j = 0; j < NUM_ACTIONS; j++) {
      const a = actionGenerator.get();
      trace.actions.push(a);
      doAction(model, a);
      trace.consequences.push(model.snapshot());
      // assert.ok(stakingWithoutSlashing(blocks), 'stakingWithoutSlashing');
      // assert.ok(
      // bondBasedConsumerVotingPower(blocks),
      // 'bondBasedConsumerVotingPower',
      // );
    }
    allEvents.push(...events);
    trace.dump(`${DIR}trace_${i}.json`);
    ////////////////////////
    elapsedMillis += end.rounded();
    if (i % 4000 === 0) {
      console.log(
        `done ${i}, traces per second ${i / (elapsedMillis / 1000)}`,
      );
    }
  }
  console.log(`ran ${Math.floor(outerEnd.seconds() / 60)} mins`);
  writeEventData(allEvents, 'Gen');
}

function fromTraces(
  fn,
  compare: undefined | ((model: Model, con) => void),
) {
  const raw = fs.readFileSync(`${fn}.json`, 'utf-8');
  const json = JSON.parse(raw);

  writeEventData(json.map((trace) => trace.events).flat(), 'PyCovering');

  const allEvents = [];
  json.forEach(({ transitions }) => {
    const events = [];
    const model = new Model(new Blocks(), events);
    transitions.forEach(([a, con]) => {
      switch (a.kind) {
        case 'Delegate':
          model.delegate(a.val, a.amt);
          break;
        case 'Undelegate':
          model.undelegate(a.val, a.amt);
          break;
        case 'JumpNBlocks':
          model.jumpNBlocks(a.n, a.chains, a.seconds_per_block);
          break;
        case 'Deliver':
          model.deliver(a.chain);
          break;
        case 'ProviderSlash':
          model.providerSlash(a.val, a.height, a.power, a.factor);
          break;
        case 'ConsumerSlash':
          model.consumerSlash(a.val, a.power, a.height, a.is_downtime);
          break;
      }
      if (compare) {
        compare(model, con);
      }
    });
    allEvents.push(...events);
  });
  writeEventData(allEvents, 'TsCovering');
}

export { gen, fromTraces };

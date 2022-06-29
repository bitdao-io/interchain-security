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
import { Events } from './events.js';

function mkdirIfNotExist(dir) {
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir);
  }
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
  height: number;
  factor: number;
};
type ConsumerSlash = {
  kind: string;
  val: number;
  power: number;
  height: number;
  isDowntime: number;
};

function weightedRandomKey(distr) {
  const x =
    Math.random() * _.reduce(_.values(distr), (sum, y) => sum + y, 0);
  const pairs = _.pairs(distr);
  let i = 0;
  let cum = 0;
  while (i < pairs.length - 1 || cum + pairs[i][1] < x) {
    i += 1;
    cum += pairs[i][1];
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

  selectUnelegate = (a): Undelegate => {
    this.undelegatedSinceBlock[a.val] = true;
    return { ...a, amt: _.random(1, 5) * TOKEN_SCALAR };
  };

  selectJumpNBlocks = (a): JumpNBlocks => {
    a = {
      ...a,
      chains: [], // TODO:
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
    // TODO:
    this.jailed[a.val] = true;
    return {
      ...a,
      power: _.random(1, 6) * TOKEN_SCALAR,
      height: Math.floor(Math.random() * this.model.h[P]),
      factor: _.sample([SLASH_DOUBLESIGN, SLASH_DOWNTIME]),
    };
  };
  selectConsumerSlash = (a): ConsumerSlash => {
    // TODO:
    this.jailed[a.val] = true;
    return {
      ...a,
      power: _.random(1, 6) * TOKEN_SCALAR,
      height: Math.floor(Math.random() * this.model.h[C]),
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
    const json = JSON.stringify({ transitions });
    fs.writeFile(fn, json, 'utf8', (err) => {
      if (err) throw err;
    });
  };
}

function doAction(model, a: Action) {
  if (a.kind === 'Delegate') {
    model.delegate(a);
  }
  if (a.kind === 'Undelegate') {
    model.undelegate(a);
  }
  if (a.kind === 'JumpNBlocks') {
    model.jumpNBlocks(a);
  }
  if (a.kind === 'Deliver') {
    model.deliver(a);
  }
  if (a.kind === 'ProviderSlash') {
    model.providerSlash(a);
  }
  if (a.kind === 'ConsumerSlash') {
    model.consumerSlash(a);
  }
}

function gen() {
  const GOAL_TIME_MINS = 5;
  const goalTimeMillis = GOAL_TIME_MINS * 60 * 1000;
  const NUM_ACTIONS = 40;
  const DIR = 'traces/';
  mkdirIfNotExist(DIR);
  let numRuns = 1000;
  let elapsed = 0;
  let i = 0;
  while (i < numRuns) {
    i += 1;
    numRuns = Math.round(goalTimeMillis / (elapsed / i) + 0.5);
    const end = timeSpan();
    ////////////////////////
    const blocks = new Blocks();
    const events = new Events();
    const model = new Model(blocks, events);
    const actionGenerator = new ActionGenerator(model);
    const trace = new Trace();
    for (let j = 0; j < NUM_ACTIONS; j++) {
      const a = actionGenerator.get();
      trace.actions.push(a);
      doAction(model, a);
      trace.consequences.push(model.snapshot());
      // todo, check properties
      stakingWithoutSlashing(blocks);
      bondBasedConsumerVotingPower(blocks);
      trace.dump(`${DIR}trace_${i}.json`);
    }
    ////////////////////////
    elapsed = end.rounded();
  }
}

export { gen };

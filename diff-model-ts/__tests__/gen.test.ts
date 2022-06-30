import { Model, Undelegation, Unval } from '../src/model.js';
import * as fs from 'fs';
import { gen } from '../src/gen.js';

function undel(o): Undelegation {
  return {
    val: o.val,
    creationHeight: o.creation_height,
    completionTime: o.completion_time,
    balance: o.balance,
    initialBalance: o.initial_balance,
    onHold: o.on_hold,
    opID: o.op_id,
    expired: o.expired,
  };
}

function unval(o): Unval {
  return {
    val: o.val,
    unbondingHeight: o.unbonding_height,
    unbondingTime: o.unbonding_time,
    onHold: o.on_hold,
    opID: o.op_id,
    expired: o.expired,
  };
}

describe('against python', () => {
  const fn = 'traces_all';
  const raw = fs.readFileSync(`${fn}.json`, 'utf-8');
  let json = JSON.parse(raw);
  const offset = 0;
  json = json.slice(offset, json.length);

  it('dt', () => {
    let j = offset;
    json.forEach(({ transitions }) => {
      const events = [];
      const model = new Model(undefined, events);

      let i = 0;
      transitions.forEach(([a, con]) => {
        console.log(`trace ${j}, action ${i}, kind ${a.kind}`);
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
        expect(model.staking.undelegationQ).toStrictEqual(
          con.undelegationQ.map(undel),
        );
        expect(model.staking.validatorQ).toStrictEqual(
          con.validatorQ.map(unval),
        );
        expect(model.h).toStrictEqual(con.h);
        expect(model.t).toStrictEqual(con.t);
        expect(model.ccvC.power).toStrictEqual(
          con.power.map((e) => (Number.isInteger(e) ? e : undefined)),
        );
        expect(model.staking.tokens).toStrictEqual(con.tokens);
        expect(model.staking.delegatorTokens).toStrictEqual(
          con.delegator_tokens,
        );
        expect(model.staking.status).toStrictEqual(con.status);
        i += 1;
      });
      j += 1;
    });
  });
});

describe('gen', () => {
  it('dt', () => {
    gen();
    expect(true).toBeTruthy();
  });
});

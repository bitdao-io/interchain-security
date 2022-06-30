import { Model, Undelegation, Unval } from '../src/model.js';
import { fromTraces, gen } from '../src/gen.js';

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

describe('compare with python', () => {
  const compare = (model: Model, con) => {
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
  };
  fromTraces('traces_covering', compare);
});

describe('gen', () => {
  it('dt', () => {
    gen();
    expect(true).toBeTruthy();
  });
});

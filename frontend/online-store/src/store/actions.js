import * as types from './mutation-types';
// 提交mutation
function makeAction (type) {
  return ({ commit }, ...args) => commit(type, ...args);
};
/**

此课程提供者：微信imax882

+微信imax882
办理会员 课程全部免费看

课程清单：https://leaaiv.cn

全网最全 最专业的 一手课程

成立十周年 会员特惠 速来抢购

**/
export const setInfo = makeAction(types.SET_INFO);
export const setNav = makeAction(types.SET_NAV);

export const setShopList = makeAction(types.SET_SHOPLIST);
// export const getPermit = makeAction(types.GET_PERMIT);
// export const getNavInfo = makeAction(types.GET_NAV);

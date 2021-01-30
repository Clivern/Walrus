/** @format */

import { setupAction, authAction, fetchInfo } from "@/common/auth.api";

const state = () => ({
	authResult: {},
	setupResult: {},
	towerInfo: {},
});

const getters = {
	getAuthResult: (state) => {
		return state.authResult;
	},
	getSetupResult: (state) => {
		return state.setupResult;
	},
	getTowerInfo: (state) => {
		return state.towerInfo;
	},
};

const actions = {
	async authAction(context, payload) {
		const result = await authAction(payload);
		context.commit("SET_AUTH_RESULT", result.data);
		return result;
	},

	async setupAction(context, payload) {
		const result = await setupAction(payload);
		context.commit("SET_SETUP_RESULT", result.data);
		return result;
	},

	async fetchInfo(context) {
		const result = await fetchInfo();
		context.commit("SET_TOWER_INFO", result.data);
		return result;
	},
};

const mutations = {
	SET_AUTH_RESULT(state, authResult) {
		state.authResult = authResult;
	},
	SET_SETUP_RESULT(state, setupResult) {
		state.setupResult = setupResult;
	},
	SET_TOWER_INFO(state, towerInfo) {
		state.towerInfo = towerInfo;
	},
};

export default {
	namespaced: true,
	state,
	getters,
	actions,
	mutations,
};

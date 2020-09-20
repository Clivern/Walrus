/** @format */

import { getTowerReadiness, getTowerHealth } from "@/common/tower.api";

const state = () => ({
	readiness: {},
	health: {},
});

const getters = {
	getTowerReadiness: (state) => {
		return state.readiness;
	},
	getTowerHealth: (state) => {
		return state.health;
	},
};

const actions = {
	async fetchTowerReadiness({ commit }) {
		const result = await getTowerReadiness();
		commit("SET_TOWER_READINESS", result.data);
		return result;
	},

	async fetchTowerHealth({ commit }) {
		const result = await getTowerHealth();
		commit("SET_TOWER_HEALTH", result.data);
		return result;
	},
};

const mutations = {
	SET_TOWER_READINESS(state, readiness) {
		state.readiness = readiness;
	},
	SET_TOWER_HEALTH(state, health) {
		state.health = health;
	},
};

export default {
	namespaced: true,
	state,
	getters,
	actions,
	mutations,
};

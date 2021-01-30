/** @format */

import { updateSettings, getSettings } from "@/common/settings.api";

const state = () => ({
	getSettingsResult: {},
	updateSettingsResult: {},
});

const getters = {
	getSettingsResult: (state) => {
		return state.getSettingsResult;
	},
	getUpdateSettingsResult: (state) => {
		return state.updateSettingsResult;
	},
};

const actions = {
	async getSettingsAction(context) {
		const result = await getSettings();
		context.commit("SET_GET_SETTINGS_RESULT", result.data);
		return result;
	},

	async updateSettingsAction(context, payload) {
		const result = await updateSettings(payload);
		context.commit("SET_UPDATE_SETTINGS_RESULT", result.data);
		return result;
	},
};

const mutations = {
	SET_GET_SETTINGS_RESULT(state, getSettingsResult) {
		state.getSettingsResult = getSettingsResult;
	},
	SET_UPDATE_SETTINGS_RESULT(state, updateSettingsResult) {
		state.updateSettingsResult = updateSettingsResult;
	},
};

export default {
	namespaced: true,
	state,
	getters,
	actions,
	mutations,
};

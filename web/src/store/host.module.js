/** @format */

import {
	getHosts,
	getHost,
	getHostCrons,
	getHostCron,
	getHostJobs,
	getHostJob,
	getAgents,
	getAgent,
	deleteHost,
	deleteHostCron,
	deleteHostJob,
	deleteAgent,
	createHostCron,
	updateHostCron,
	updateHostJob,
} from "@/common/host.api";

const state = () => ({
	getHostsResult: {},
	getHostResult: {},
	getHostCronsResult: {},
	getHostCronResult: {},
	getHostJobsResult: {},
	getHostJobResult: {},
	getAgentsResult: {},
	getAgentResult: {},
	deleteHostResult: {},
	deleteHostCronResult: {},
	deleteHostJobResult: {},
	deleteAgentResult: {},
	createHostCronResult: {},
	updateHostCronResult: {},
	updateHostJobResult: {},
});

const getters = {
	getHostsResult: (state) => {
		return state.getHostsResult;
	},
	getHostResult: (state) => {
		return state.getHostResult;
	},
	getHostCronsResult: (state) => {
		return state.getHostCronsResult;
	},
	getHostCronResult: (state) => {
		return state.getHostCronResult;
	},
	getHostJobsResult: (state) => {
		return state.getHostJobsResult;
	},
	getHostJobResult: (state) => {
		return state.getHostJobResult;
	},
	getAgentsResult: (state) => {
		return state.getAgentsResult;
	},
	getAgentResult: (state) => {
		return state.getAgentResult;
	},
	getDeleteHostResult: (state) => {
		return state.deleteHostResult;
	},
	getDeleteHostCronResult: (state) => {
		return state.deleteHostCronResult;
	},
	getDeleteHostJobResult: (state) => {
		return state.deleteHostJobResult;
	},
	getDeleteAgentResult: (state) => {
		return state.deleteAgentResult;
	},
	getCreateHostCronResult: (state) => {
		return state.createHostCronResult;
	},
	getUpdateHostCronResult: (state) => {
		return state.updateHostCronResult;
	},
	getUpdateHostJobResult: (state) => {
		return state.updateHostJobResult;
	},
};

const actions = {
	async getHostsAction(context) {
		const result = await getHosts();
		context.commit("SET_GET_HOSTS_RESULT", result.data);
		return result;
	},
	async getHostAction(context, payload) {
		const result = await getHost(payload);
		context.commit("SET_GET_HOST_RESULT", result.data);
		return result;
	},
	async getHostCronsAction(context, payload) {
		const result = await getHostCrons(payload);
		context.commit("SET_GET_HOST_CRONS_RESULT", result.data);
		return result;
	},
	async getHostCronAction(context, payload) {
		const result = await getHostCron(payload);
		context.commit("SET_GET_HOST_CRON_RESULT", result.data);
		return result;
	},
	async getHostJobsAction(context, payload) {
		const result = await getHostJobs(payload);
		context.commit("SET_GET_HOST_JOBS_RESULT", result.data);
		return result;
	},
	async getHostJobAction(context, payload) {
		const result = await getHostJob(payload);
		context.commit("SET_GET_HOST_JOB_RESULT", result.data);
		return result;
	},
	async getAgentsAction(context, payload) {
		const result = await getAgents(payload);
		context.commit("SET_GET_AGENTS_RESULT", result.data);
		return result;
	},
	async getAgentAction(context, payload) {
		const result = await getAgent(payload);
		context.commit("SET_GET_AGENT_RESULT", result.data);
		return result;
	},
	async deleteHostAction(context, payload) {
		const result = await deleteHost(payload);
		context.commit("SET_DELETE_HOST_RESULT", result.data);
		return result;
	},
	async deleteHostCronAction(context, payload) {
		const result = await deleteHostCron(payload);
		context.commit("SET_DELETE_HOST_CRON_RESULT", result.data);
		return result;
	},
	async deleteHostJobAction(context, payload) {
		const result = await deleteHostJob(payload);
		context.commit("SET_DELETE_HOST_JOB_RESULT", result.data);
		return result;
	},
	async deleteAgentAction(context, payload) {
		const result = await deleteAgent(payload);
		context.commit("SET_DELETE_AGENT_RESULT", result.data);
		return result;
	},
	async createHostCronAction(context, payload) {
		const result = await createHostCron(payload);
		context.commit("SET_CREATE_HOST_CRON_RESULT", result.data);
		return result;
	},
	async updateHostCronAction(context, payload) {
		const result = await updateHostCron(payload);
		context.commit("SET_UPDATE_HOST_CRON_RESULT", result.data);
		return result;
	},
	async updateHostJobAction(context, payload) {
		const result = await updateHostJob(payload);
		context.commit("SET_UPDATE_HOST_JOB_RESULT", result.data);
		return result;
	},
};

const mutations = {
	SET_GET_HOSTS_RESULT(state, getHostsResult) {
		state.getHostsResult = getHostsResult;
	},
	SET_GET_HOST_RESULT(state, getHostResult) {
		state.getHostResult = getHostResult;
	},
	SET_GET_HOST_CRONS_RESULT(state, getHostCronsResult) {
		state.getHostCronsResult = getHostCronsResult;
	},
	SET_GET_HOST_CRON_RESULT(state, getHostCronResult) {
		state.getHostCronResult = getHostCronResult;
	},
	SET_GET_HOST_JOBS_RESULT(state, getHostJobsResult) {
		state.getHostJobsResult = getHostJobsResult;
	},
	SET_GET_HOST_JOB_RESULT(state, getHostJobResult) {
		state.getHostJobResult = getHostJobResult;
	},
	SET_GET_AGENTS_RESULT(state, getAgentsResult) {
		state.getAgentsResult = getAgentsResult;
	},
	SET_GET_AGENT_RESULT(state, getAgentResult) {
		state.getAgentResult = getAgentResult;
	},
	SET_DELETE_HOST_RESULT(state, deleteHostResult) {
		state.deleteHostResult = deleteHostResult;
	},
	SET_DELETE_HOST_CRON_RESULT(state, deleteHostCronResult) {
		state.deleteHostCronResult = deleteHostCronResult;
	},
	SET_DELETE_HOST_JOB_RESULT(state, deleteHostJobResult) {
		state.deleteHostJobResult = deleteHostJobResult;
	},
	SET_DELETE_AGENT_RESULT(state, deleteAgentResult) {
		state.deleteAgentResult = deleteAgentResult;
	},
	SET_CREATE_HOST_CRON_RESULT(state, createHostCronResult) {
		state.createHostCronResult = createHostCronResult;
	},
	SET_UPDATE_HOST_CRON_RESULT(state, updateHostCronResult) {
		state.updateHostCronResult = updateHostCronResult;
	},
	SET_UPDATE_HOST_JOB_RESULT(state, updateHostJobResult) {
		state.updateHostJobResult = updateHostJobResult;
	},
};

export default {
	namespaced: true,
	state,
	getters,
	actions,
	mutations,
};

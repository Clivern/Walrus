/** @format */

import ApiService from "./api.service.js";

const getHosts = () => {
	return ApiService.get("/api/v1/host");
};

const getHost = (payload) => {
	return ApiService.get("/api/v1/host/" + payload["hostname"]);
};

const deleteHost = (payload) => {
	return ApiService.delete("/api/v1/host/" + payload["hostname"]);
};

const createHostCron = (payload) => {
	return ApiService.post(
		"/api/v1/host/" + payload["hostname"] + "/cron",
		payload
	);
};

const getHostCrons = (payload) => {
	return ApiService.get("/api/v1/host/" + payload["hostname"] + "/cron");
};

const getHostCron = (payload) => {
	return ApiService.get(
		"/api/v1/host/" + payload["hostname"] + "/cron/" + payload["cronId"]
	);
};

const updateHostCron = (payload) => {
	return ApiService.put(
		"/api/v1/host/" + payload["hostname"] + "/cron/" + payload["cronId"],
		payload
	);
};

const deleteHostCron = (payload) => {
	return ApiService.delete(
		"/api/v1/host/" + payload["hostname"] + "/cron/" + payload["cronId"]
	);
};

const getHostJobs = (payload) => {
	return ApiService.get("/api/v1/host/" + payload["hostname"] + "/job");
};

const getHostJob = (payload) => {
	return ApiService.get(
		"/api/v1/host/" + payload["hostname"] + "/job/" + payload["jobId"]
	);
};

const updateHostJob = (payload) => {
	return ApiService.put(
		"/api/v1/host/" + payload["hostname"] + "/job/" + payload["jobId"],
		payload
	);
};

const deleteHostJob = (payload) => {
	return ApiService.delete(
		"/api/v1/host/" + payload["hostname"] + "/job/" + payload["jobId"]
	);
};

const getAgents = (payload) => {
	return ApiService.get("/api/v1/host/" + payload["hostname"] + "/agent");
};

const getAgent = (payload) => {
	return ApiService.get(
		"/api/v1/host/" + payload["hostname"] + "/agent/" + payload["agentId"]
	);
};

const deleteAgent = (payload) => {
	return ApiService.delete(
		"/api/v1/host/" + payload["hostname"] + "/agent/" + payload["agentId"]
	);
};

export {
	getHosts,
	getHost,
	deleteHost,
	createHostCron,
	getHostCrons,
	getHostCron,
	updateHostCron,
	deleteHostCron,
	getHostJobs,
	getHostJob,
	updateHostJob,
	deleteHostJob,
	getAgents,
	getAgent,
	deleteAgent,
};

/** @format */

import ApiService from "./api.service.js";

const setupAction = (payload) => {
	return ApiService.post("/action/setup", payload);
};

const authAction = (payload) => {
	return ApiService.post("/action/auth", payload);
};

const fetchInfo = () => {
	return ApiService.get("/action/info");
};

export { setupAction, authAction, fetchInfo };

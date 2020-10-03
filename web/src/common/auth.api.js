/** @format */

import ApiService from "./api.service.js";

const setupAction = (email, password) => {
	return ApiService.post("/action/setup", {
		email: email,
		password: password,
	});
};

const authAction = (email, password) => {
	return ApiService.post("/action/auth", {
		email: email,
		password: password,
	});
};

const fetchInfo = () => {
	return ApiService.get("/action/info");
};

export { setupAction, authAction, fetchInfo };

/** @format */

import ApiService from "./api.service.js";

const setupAction = (email, password) => {
	return ApiService.post("/setup", {
		email: email,
		password: password,
	});
};

const authAction = (email, password) => {
	return ApiService.post("/auth", {
		email: email,
		password: password,
	});
};

const fetchInfo = () => {
	return ApiService.get("/info");
};

export { setupAction, authAction, fetchInfo };

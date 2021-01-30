/** @format */

import ApiService from "./api.service.js";

const updateSettings = (payload) => {
	return ApiService.put("/api/v1/settings", payload);
};

const getSettings = () => {
	return ApiService.get("/api/v1/settings");
};

export { updateSettings, getSettings };

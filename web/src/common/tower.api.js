/** @format */

import ApiService from "./api.service.js";

const getTowerReadiness = () => {
	return ApiService.get("/_ready");
};

const getTowerHealth = () => {
	return ApiService.get("/_health");
};

export { getTowerReadiness, getTowerHealth };

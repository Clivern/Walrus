/** @format */

import ApiService from "./api.service.js";

const getItem = () => {
	return ApiService.get("/");
};

export { getItem };

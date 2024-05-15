import Axios from "axios";

const headers = {
  "Content-Type": "application/json",
};

const baseURL = `${import.meta.env.VITE_BACKEND_URL}/api`;

export const axiosPublic = Axios.create({
  baseURL,
  headers,
});
export const axiosPrivate = Axios.create({
  baseURL,
  headers,
  withCredentials: true,
});

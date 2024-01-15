import axios from "axios";

const api = axios.create({
  baseURL: process.env.REACT_APP_SERVER_URL,
  withCredentials: true,
});

export async function makeRequest(url, options) {
  try {
    const res = await api.get(url, options);
    return res.data;
  } catch (error) {
    throw error;
  }
}

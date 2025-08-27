import api from "./api"

export const fetchDailyTrends = async () => {
  const response = await api.get("/trend/daily")
  return response.data
}

export const fetchWeeklyTrends = async () => {
  const response = await api.get("/trend/weekly")
  return response.data
}
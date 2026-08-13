const klookCityIds: Record<string, string> = {
  台北: "19",
  新北: "6488",
  桃園: "4737",
  花蓮: "20",
  高雄: "22",
  台中: "25",
  台南: "164",
  宜蘭: "42",
  澎湖: "43",
  台東: "47",
  嘉義: "436",
  屏東: "7992",
  苗栗: "8109",
  基隆: "10048",
  雲林: "16222",
  彰化: "24415",
  南投: "25303",
};

const hsinchuCityTownships = new Set(["東區", "北區", "香山區"]);

export const resolveKlookCityId = (cityName?: string, townshipName?: string) => {
  if (!cityName) return "";
  if (cityName !== "新竹") return klookCityIds[cityName] || "";
  return townshipName && hsinchuCityTownships.has(townshipName)
    ? "27456"
    : "17312";
};

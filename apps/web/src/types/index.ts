export type Draw = {
  id: number;
  drawDate: string;
  n1: number;
  n2: number;
  n3: number;
  n4: number;
  n5: number;
  star1: number;
  star2: number;
  createdAt: string;
  updatedAt: string;
};

export type DrawPayload = {
  drawDate: string;
  n1: number;
  n2: number;
  n3: number;
  n4: number;
  n5: number;
  star1: number;
  star2: number;
};

export type PaginatedDraws = {
  draws: Draw[];
  total: number;
  page: number;
  limit: number;
};

export type NumberFrequency = {
  value: number;
  count: number;
};

export type FrequenciesResponse = {
  numbers: NumberFrequency[];
  stars: NumberFrequency[];
};

export type PositionFrequency = {
  position: string;
  values: NumberFrequency[];
};

export type PositionsResponse = {
  numbers: PositionFrequency[];
  stars: PositionFrequency[];
};

export type DelayStat = {
  value: number;
  delay: number;
};

export type DelaysResponse = {
  numbers: DelayStat[];
  stars: DelayStat[];
};

export type PairStat = {
  a: number;
  b: number;
  count: number;
};

export type HotColdResponse = {
  hotNumbers: NumberFrequency[];
  coldNumbers: NumberFrequency[];
  hotStars: NumberFrequency[];
  coldStars: NumberFrequency[];
};

export type DashboardStats = {
  totalDraws: number;
  mostFrequentNumber?: NumberFrequency;
  mostFrequentStar?: NumberFrequency;
  mostDelayedNumber?: DelayStat;
  lastDrawDate: string;
};

export type GenerationStrategy = "hot" | "cold" | "delayed" | "balanced" | "random" | "anti_human";

export type GeneratedCombination = {
  numbers: number[];
  stars: number[];
  strategy: GenerationStrategy;
  score: number;
  explanation: string;
};

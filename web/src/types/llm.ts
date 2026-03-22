export type LLMModelsConfig = {
  default: string;
  sonnet: string;
  opus: string;
  haiku: string;
};

export type LLMConfig = {
  apiKey: string;
  baseUrl: string;
  models: LLMModelsConfig;
};

export type LLMProvider = {
  id: string;
  name: string;
  baseUrl: string;
  models: string[];
};

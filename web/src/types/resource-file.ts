export type ResourceKind = "agents" | "commands" | "rules";

export type ResourceListResponse = {
  items: string[];
};

export type ResourceFile = {
  name: string;
  content: string;
};

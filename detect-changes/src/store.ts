export type Deployment = {
  key: string;
  created: Date;
  commitSHA: string;
};

export interface Store {
  get(key: string, altKeys: string[]): Promise<Deployment>;
  set(deploy: Deployment): Promise<void>;
}

import * as cache from '@actions/cache';
import * as core from '@actions/core';
import * as fs from 'node:fs/promises';
import * as path from 'node:path';
import * as os from 'node:os';
import { Deployment, Store } from './store';

export class GHCacheStore implements Store {
  private readonly _cachepath: string;
  private readonly _storepath: string;

  constructor() {
    this._cachepath = path.join(os.homedir(), '.cache/snapaddy-actions-detect-changes');
    this._storepath = path.join(this._cachepath, 'store.json');
  }

  async get(key: string): Promise<Deployment> {
    await cache.restoreCache([this._cachepath], key, []);

    core.debug(`Trying to read cachefile. file=${this._storepath}`);
    const content = await fs.readFile(this._storepath, { encoding: 'utf-8' });

    core.debug(`Trying to parse cachefile content. content=${content}`);
    const dpl: Deployment = JSON.parse(content);

    if (dpl.key !== key) {
      throw new Error(`Deployment does not match provided key. got=${dpl.key} want=${key}`);
    }
    core.debug(`Got Deployment.`);

    return dpl;
  }

  async set(deploy: Deployment): Promise<void> {
    core.debug(`Trying to persist deployment to cachefile. file=${this._storepath}`);

    await fs.mkdir(this._cachepath, { recursive: true });
    await fs.writeFile(this._storepath, JSON.stringify(deploy), { encoding: 'utf8' });
    core.debug(`Persisted Deployment.`);

    await cache.saveCache([this._cachepath], deploy.key);
  }
}

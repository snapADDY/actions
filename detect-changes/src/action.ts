import * as core from '@actions/core';
import * as github from '@actions/github';
import * as cp from 'node:child_process';
import { Inputs, KEY_PREFIX, Outputs } from './constants';
import { Deployment, Store } from './store';
import { GHCacheStore } from './store-ghcache';

export async function run(): Promise<void> {
  const changeKey: string = KEY_PREFIX + core.getInput(Inputs.ChangeKey);

  try {
    const store: Store = new GHCacheStore();
    let dpl = await store.get(changeKey);

    core.setOutput(Outputs.LastCommitSHA, dpl.commitSHA);

    let files = cp.execSync(`git diff --name-only ${dpl.commitSHA}`).toString();
    core.setOutput(Outputs.ChangedFiles, files);
  } catch (err) {
    core.info(`Could not find any deployment. key=${changeKey} error=${(err as Error).message}`);

    core.setOutput(Outputs.ChangedFiles, '**');
    core.setOutput(Outputs.LastCommitSHA, '');
  }
}

export async function postRun(): Promise<void> {
  if (core.getInput(Inputs.PersistRun) !== 'true') {
    core.info('persist-run is disabled');

    return;
  }

  const context = github.context;

  const dpl: Deployment = {
    key: KEY_PREFIX + core.getInput(Inputs.ChangeKey),
    created: new Date(),
    commitSHA: context.sha,
  };

  core.debug(`Created Deployment. deployment=${JSON.stringify(dpl)}`);

  const store: Store = new GHCacheStore();
  await store.set(dpl);

  core.info(`Successfully persisted deployment. key=${dpl.key} sha=${dpl.commitSHA}`);
}

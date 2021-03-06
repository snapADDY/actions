import * as core from '@actions/core';
import * as github from '@actions/github';
import * as cp from 'node:child_process';
import { Inputs, KEY_PREFIX, Outputs } from './constants';
import { Deployment, Store } from './store';
import { GHCacheStore } from './store-ghcache';

export async function run(): Promise<void> {
  const changeKeyBase: string = KEY_PREFIX + core.getInput(Inputs.ChangeKey) + '-';
  const changeKey = changeKeyBase + github.context.sha;
  try {
    const store: Store = new GHCacheStore();
    let dpl = await store.get(changeKey, [changeKeyBase]);

    core.info(`Found previous deployment. key=${changeKey} sha=${dpl.commitSHA}`);
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

  try {
    const context = github.context;

    const changeKey = KEY_PREFIX + core.getInput(Inputs.ChangeKey) + '-' + context.sha;
    const dpl: Deployment = {
      key: changeKey,
      created: new Date(),
      commitSHA: context.sha,
    };
    core.debug(`Created Deployment. deployment=${JSON.stringify(dpl)}`);

    const store: Store = new GHCacheStore();
    await store.set(dpl);
    core.info(`Successfully persisted deployment. key=${dpl.key} sha=${dpl.commitSHA}`);
  } catch (err) {
    core.info(`Could persist deployment. error=${(err as Error).message}`);
  }
}

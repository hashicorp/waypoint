/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { TestContext } from 'ember-test-helpers';
import { render, settled } from '@ember/test-helpers';
import { GetJobStreamRequest, GetJobStreamResponse } from 'waypoint-pb';
import { getTerminalText } from '../../helpers/xterm';
import Service from '@ember/service';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | operation-logs', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(async function (this: TestContext) {
    this.owner.register('service:api', ApiStub);
  });

  test('happy path', async function (assert) {
    let api: ApiStub = this.owner.lookup('service:api');

    await render(hbs`
      <OperationLogs @jobId="test-job-1" />
    `);

    api.sendLine('Test message');
    await settled();

    let text = await getTerminalText();

    assert.equal(text, 'Test message');
  });

  test('changing @jobId', async function (assert) {
    let api: ApiStub = this.owner.lookup('service:api');

    // Render the component with test-job-1

    this.set('jobId', 'test-job-1');
    await render(hbs`
      <OperationLogs @jobId={{this.jobId}} />
    `);
    api.sendLine('Message from test-job-1');
    await settled();

    let firstStream = api.client.stream;

    // Send in a new job ID

    this.set('jobId', 'test-job-2');
    await settled();

    api.sendLine('Message from test-job-2');
    await settled();

    let text = await getTerminalText();

    assert.equal(api.currentRequest?.getJobId(), 'test-job-2', 'We requested the job stream for the new ID');
    assert.equal(text, 'Message from test-job-2', 'We cleared the terminal');
    assert.ok(firstStream?.cancelled, 'We cancelled the first stream');
  });
});

// Stubs

class ApiStub extends Service {
  client = new ApiClientStub();

  WithMeta() {
    return {};
  }

  sendLine(message: string) {
    this.client.sendLine(message);
  }

  get currentRequest() {
    return this.client.currentRequest;
  }
}

class ApiClientStub {
  stream?: StreamStub;
  currentRequest?: GetJobStreamRequest;

  getJobStream(request: GetJobStreamRequest) {
    this.currentRequest = request;
    this.stream = new StreamStub();

    return this.stream;
  }

  sendLine(message: string) {
    this.stream?.sendLine(message);
  }
}

type StreamEvent = 'data' | 'status';
type StreamHandler = (response: GetJobStreamResponse) => void;

class StreamStub {
  cancelled = false;

  handlers: Record<StreamEvent, StreamHandler[]> = {
    data: [],
    status: [],
  };

  on(event: StreamEvent, handler: StreamHandler): void {
    this.handlers[event].push(handler);
  }

  cancel() {
    this.cancelled = true;
  }

  sendLine(message: string) {
    let response = new GetJobStreamResponse();
    let terminal = new GetJobStreamResponse.Terminal();
    let event = new GetJobStreamResponse.Terminal.Event();
    let line = new GetJobStreamResponse.Terminal.Event.Line();

    response.setTerminal(terminal);
    terminal.setEventsList([event]);
    event.setLine(line);
    line.setMsg(message);

    for (let handler of this.handlers.data) {
      handler(response);
    }
  }
}

const worker = new Worker('../worker/worker.js', { type: 'module' });

const resultView = document.querySelector('#result-view');
const errorBanner = document.querySelector('#error-banner');
const eventLog = document.querySelector('#event-log');

let nextRequestId = 1;
const pending = new Map();

worker.addEventListener('message', (messageEvent) => {
  const msg = messageEvent.data;
  if (!msg || typeof msg !== 'object') {
    return;
  }

  if (msg.type === 'event') {
    appendEvent(msg.event, msg.payload);
    return;
  }

  if (msg.type === 'response') {
    const request = pending.get(msg.id);
    if (!request) return;
    pending.delete(msg.id);

    if (msg.error) {
      setError(msg.error.message || 'Unknown error');
      request.reject(new Error(msg.error.message || 'Unknown error'));
      return;
    }

    clearError();
    request.resolve(msg.result);
  }
});

initialize();
bindCoreFlowForms();
bindCommandForm();

async function initialize() {
  try {
    await request('worker.init', {
      wasmUrl: '../worker/ipatool.wasm',
    });
    appendEvent('worker.ready', { ok: true });
  } catch (error) {
    setError(`Worker init failed: ${error.message}`);
  }
}

function bindCoreFlowForms() {
  document.querySelectorAll('.controls form').forEach((form) => {
    form.addEventListener('submit', async (event) => {
      event.preventDefault();
      const command = form.dataset.command;
      const payload = collectFormPayload(form);

      try {
        const result = await request(command, payload);
        renderResult(result);

        if (command === 'download.ipa') {
          saveDownloadBlob(result);
        }
      } catch {
        // Error state already rendered.
      }
    });
  });
}

function bindCommandForm() {
  const commandForm = document.querySelector('#command-form');
  commandForm.addEventListener('submit', async (event) => {
    event.preventDefault();

    const command = document.querySelector('#command-name').value.trim();
    const payloadInput = document.querySelector('#command-payload').value;

    try {
      const payload = JSON.parse(payloadInput);
      const result = await request(command, payload);
      renderResult(result);
    } catch (error) {
      setError(error.message);
    }
  });
}

function collectFormPayload(form) {
  const data = new FormData(form);
  const payload = {};
  for (const [key, value] of data.entries()) {
    if (value === '') continue;
    const parsed = Number(value);
    payload[key] = Number.isFinite(parsed) && String(parsed) === String(value) ? parsed : value;
  }

  form.querySelectorAll('[data-sensitive="true"]').forEach((input) => {
    input.autocomplete = 'off';
  });

  return payload;
}

function request(command, payload) {
  const id = `req-${nextRequestId++}`;

  return new Promise((resolve, reject) => {
    pending.set(id, { resolve, reject });
    worker.postMessage({
      type: 'request',
      id,
      command,
      payload,
    });
  });
}

function renderResult(result) {
  resultView.textContent = JSON.stringify(result, null, 2);
}

function appendEvent(eventName, payload) {
  const li = document.createElement('li');
  li.textContent = `${new Date().toISOString()} :: ${eventName} :: ${JSON.stringify(payload)}`;
  eventLog.prepend(li);
}

function setError(message) {
  errorBanner.textContent = message;
}

function clearError() {
  errorBanner.textContent = '';
}

async function saveDownloadBlob(result) {
  const bytes = result?.fileBytes;
  if (!bytes || !Array.isArray(bytes)) {
    return;
  }

  const stream = new ReadableStream({
    start(controller) {
      controller.enqueue(new Uint8Array(bytes));
      controller.close();
    },
  });

  const response = new Response(stream);
  const blob = await response.blob();
  const href = URL.createObjectURL(blob);

  const anchor = document.createElement('a');
  anchor.href = href;
  anchor.download = result.fileName || 'download.ipa';
  anchor.click();

  setTimeout(() => URL.revokeObjectURL(href), 2500);
}

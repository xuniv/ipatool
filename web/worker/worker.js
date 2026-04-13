const state = {
  wasm: null,
  exports: null,
};

self.addEventListener('message', async (event) => {
  const msg = event.data;
  if (!msg || msg.type !== 'request') {
    return;
  }

  const { id, command, payload } = msg;
  emit('request.received', { id, command });

  try {
    const result = await handleCommand(command, payload);
    self.postMessage({
      type: 'response',
      id,
      result,
    });
    emit('request.completed', { id, command });
  } catch (error) {
    emit('request.failed', { id, command, error: error.message });
    self.postMessage({
      type: 'response',
      id,
      error: { message: error.message || 'Worker command failed' },
    });
  }
});

async function handleCommand(command, payload) {
  switch (command) {
    case 'worker.init':
      return initWasm(payload?.wasmUrl);
    case 'auth.login':
      emit('auth.started', { appleId: payload?.appleId });
      return bridgeProtocol('login', payload);
    case 'search.apps':
      emit('search.started', { term: payload?.term });
      return bridgeProtocol('search', payload);
    case 'download.ipa':
      emit('download.started', { appId: payload?.appId });
      return bridgeProtocol('download', payload);
    default:
      return bridgeProtocol(command, payload);
  }
}

async function initWasm(wasmUrl) {
  if (!wasmUrl) {
    throw new Error('wasmUrl is required');
  }

  try {
    const response = await fetch(wasmUrl);
    if (!response.ok) {
      throw new Error(`Failed to fetch wasm: HTTP ${response.status}`);
    }
    const bytes = await response.arrayBuffer();
    const module = await WebAssembly.compile(bytes);
    const instance = await WebAssembly.instantiate(module, {
      env: {
        abort() {
          throw new Error('WASM aborted');
        },
      },
    });

    state.wasm = module;
    state.exports = instance.exports || null;

    emit('worker.wasm.loaded', { wasmUrl });
    return { loaded: true };
  } catch (error) {
    emit('worker.wasm.load_failed', { wasmUrl, error: error.message });
    throw error;
  }
}

async function bridgeProtocol(operation, payload) {
  emit('bridge.request', { operation });

  if (operation === 'login') {
    return {
      ok: true,
      account: payload?.appleId,
      tokenInMemoryOnly: true,
    };
  }

  if (operation === 'search') {
    return {
      ok: true,
      items: [
        {
          appId: 284882215,
          name: `Sample result for ${payload?.term || 'query'}`,
        },
      ],
      limit: payload?.limit || 5,
    };
  }

  if (operation === 'download') {
    const fakeBytes = Array.from(new TextEncoder().encode('fake ipa bytes'));
    return {
      ok: true,
      fileName: payload?.fileName || 'app.ipa',
      fileBytes: fakeBytes,
    };
  }

  if (state.exports && typeof state.exports.handle === 'function') {
    emit('bridge.wasm.invoke', { operation });
    const rc = state.exports.handle();
    return { ok: rc === 0, code: rc };
  }

  return {
    ok: true,
    operation,
    payload,
    via: 'worker-bridge-fallback',
  };
}

function emit(eventName, payload = {}) {
  self.postMessage({
    type: 'event',
    event: eventName,
    payload,
  });
}

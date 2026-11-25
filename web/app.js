// UI layer: keeps track of what the user typed and calls the Go functions on
// window.calcWasm for every calculation. No math happens here.

const $ = (id) => document.getElementById(id);
const keys = $("keys");
const exprEl = $("expr");
const valueEl = $("value");

let current = "0";      // number being typed, as a string
let stored = null;      // left-hand number of a pending a ○ b
let pendingOp = null;   // { op, symbol } waiting for the right-hand number
let freshInput = true;  // next digit replaces `current` instead of appending

// Calls a WASM function. Returns the number, or null after showing the error.
function callWasm(op, ...args) {
  const fn = window.calcWasm?.[op];
  if (typeof fn !== "function") {
    showError(`"${op}" is not registered in main.go`);
    return null;
  }
  const res = fn(...args);
  if (!res.ok) {
    showError(res.error);
    return null;
  }
  return res.value;
}

function format(n) {
  if (!Number.isFinite(n)) return n > 0 ? "∞" : n < 0 ? "-∞" : "NaN";
  return String(parseFloat(n.toPrecision(12)));
}

function showError(msg) {
  valueEl.textContent = msg;
  valueEl.classList.add("error");
  current = "0";
  stored = null;
  pendingOp = null;
  freshInput = true;
  exprEl.textContent = "";
}

function render() {
  valueEl.classList.remove("error");
  valueEl.textContent = current;
  exprEl.textContent = pendingOp ? `${format(stored)} ${pendingOp.symbol}` : "";
}

function inputDigit(d) {
  if (freshInput) {
    current = d === "." ? "0." : d;
    freshInput = false;
  } else if (d === "." && current.includes(".")) {
    return;
  } else {
    current = current === "0" && d !== "." ? d : current + d;
  }
  render();
}

// Finishes a pending a ○ b. Returns false if WASM reported an error.
function resolvePending() {
  if (!pendingOp) return true;
  const value = callWasm(pendingOp.op, stored, Number(current));
  if (value === null) return false;
  current = format(value);
  stored = null;
  pendingOp = null;
  return true;
}

function handleOp(btn) {
  const { op, arity, symbol } = btn.dataset;

  if (arity === "1") {
    const value = callWasm(op, Number(current));
    if (value === null) return;
    current = format(value);
    freshInput = true;
    render();
    return;
  }

  // Two-number op: finish any earlier op first so "2 + 3 × " shows 5 ×.
  if (pendingOp && !freshInput && !resolvePending()) return;
  stored = Number(current);
  pendingOp = { op, symbol: symbol ?? op };
  freshInput = true;
  render();
}

function handleAction(action) {
  switch (action) {
    case "clear":
      current = "0"; stored = null; pendingOp = null; freshInput = true;
      break;
    case "back":
      if (freshInput) return;
      current = current.length > 1 ? current.slice(0, -1) : "0";
      if (current === "-") current = "0";
      break;
    case "negate":
      if (current !== "0") current = current.startsWith("-") ? current.slice(1) : "-" + current;
      break;
    case "equals":
      if (!resolvePending()) return;
      freshInput = true;
      break;
  }
  render();
}

keys.addEventListener("click", (e) => {
  const btn = e.target.closest("button");
  if (!btn || btn.disabled) return;
  if (btn.dataset.digit) inputDigit(btn.dataset.digit);
  else if (btn.dataset.op) handleOp(btn);
  else if (btn.dataset.action) handleAction(btn.dataset.action);
});

// Keyboard support.
const keyMap = {
  "+": '[data-op="add"]', "-": '[data-op="subtract"]',
  "*": '[data-op="multiply"]', "/": '[data-op="divide"]',
  "^": '[data-op="power"]', "%": '[data-op="percent"]',
  "Enter": '[data-action="equals"]', "=": '[data-action="equals"]',
  "Backspace": '[data-action="back"]', "Escape": '[data-action="clear"]',
};
document.addEventListener("keydown", (e) => {
  if (e.metaKey || e.ctrlKey || e.altKey) return;
  const sel = /^[0-9.]$/.test(e.key) ? `[data-digit="${e.key}"]` : keyMap[e.key];
  const btn = sel && keys.querySelector(sel);
  if (!btn) return;
  e.preventDefault();
  btn.click();
});

// Lists every data-op button and whether its Go function works yet, by
// calling it once with sample input.
function renderStatus() {
  const items = [...keys.querySelectorAll("[data-op]")].map((btn) => {
    const { op, arity } = btn.dataset;
    const fn = window.calcWasm[op];
    let state, note;
    if (typeof fn !== "function") {
      state = "missing"; note = "not registered in main.go";
    } else {
      const res = arity === "1" ? fn(4) : fn(6, 3);
      if (!res.ok && res.error === "not implemented yet") {
        state = "todo"; note = "TODO in calc.go";
      } else {
        state = "done"; note = "implemented";
      }
    }
    const li = document.createElement("li");
    li.className = state;
    const code = document.createElement("code");
    code.textContent = op;
    li.append(code, ` ${note}`);
    return li;
  });
  $("status").replaceChildren(...items);
}

// Load the WASM module. Go's main() fires "wasm-ready" once calcWasm exists.
keys.querySelectorAll("button").forEach((b) => (b.disabled = true));
window.addEventListener("wasm-ready", () => {
  keys.querySelectorAll("button").forEach((b) => (b.disabled = false));
  renderStatus();
  render();
});

const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
  .then(({ instance }) => go.run(instance))
  .catch((err) => {
    $("status").innerHTML = "";
    const li = document.createElement("li");
    li.className = "missing";
    li.textContent = `Could not load main.wasm: ${err.message}. Did you run "make build"?`;
    $("status").append(li);
  });

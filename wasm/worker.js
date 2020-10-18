addEventListener('message', async (e) => {
	// initialize the Go WASM glue
	const go = new self.Go();

	// e.data contains the code from the main thread
	const result = await WebAssembly.instantiate(fetch('main.wasm'), go.importObject);

	// hijack the console.log function to capture stdout
	let oldLog = console.log;
	// send each line of output to the main thread
	console.log = (line) => { postMessage({
		message: line
	}); };

	// run the code
	await go.run(result.instance);
	console.log = oldLog;

	// tell the main thread we are done
	postMessage({
		done: true
	});
}, false);

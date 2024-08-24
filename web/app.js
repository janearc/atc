(async () => {
    const go = new Go();
    const wasmModule = await WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject);
    go.run(wasmModule.instance);

    // Assuming Go will provide the activities via a function called getActivities
    const activities = getActivities();

    const table = document.getElementById("activities-table");

    activities.forEach(activity => {
        const row = table.insertRow();
        const typeCell = row.insertCell(0);
        const durationCell = row.insertCell(1);
        const tssCell = row.insertCell(2);

        typeCell.innerHTML = activity.type;
        durationCell.innerHTML = Math.round(activity.movingTime / 60);
        tssCell.innerHTML = activity.tss.toFixed(2);
    });
})();
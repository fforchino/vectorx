function LoadBotControlPage() {
    let botSerialNo = QueryParams.esn;
    document.getElementById("bc_serial_no").value = botSerialNo;
    document.getElementById("bc_title").value = "Play with "+botSerialNo;
}

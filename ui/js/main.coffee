FileListRequestProtocol = "fsrp"
webSocket = new WebSocket("ws://localhost:1241/ws")
fileList = null

webSocket.onmessage = (event) ->
  message = JSON.parse(event.data)
  logger = document.getElementsByName("logger")
  fileList = message.N
  dispTable(fileList)
  return

webSocket.onopen = () ->
  webSocket.send(FileListRequestProtocol)
  return

dispTable = (list) ->
  table = document.getElementById("fileList")
  i=0
  for name, file of list.Children
    row = table.insertRow(i)
    i++
    addFileToRow(row, file)
  return

addFileToRow = (row, file) ->
  name = row.insertCell(0)
  isDir = row.insertCell(1)
  atLocal = row.insertCell(2)
  size = row.insertCell(3)

  name.innerHTML = file.Name
  isDir.innerHTML = file.IsDir
  atLocal.innerHTML = file.AtLocal
  size.innerHTML = file.Size
  return

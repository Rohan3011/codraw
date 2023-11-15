let socket = null

/**
 * @type {HTMLCanvasElement}
 */
const canvas = document.getElementById('canvas')

/**
 * @type {HTMLInputElement}
 */
const colorPicker = document.getElementById('colorPicker')

/**
 * @type {HTMLButtonElement}
 */
const clearButton = document.getElementById('clearButton')

const ctx = canvas.getContext('2d')
let drawing = false
let prevPoint = null

function initSocket () {
  try {
    if (!window['WebSocket']) {
      console.error('Error: Your browser does not support web sockets.')
    }

    socket = new WebSocket(`ws://${window.location.host}/ws`)

    socket.onopen = function () {
      socket.send(
        JSON.stringify({
          type: 'client-ready',
          data: null
        })
      )
      console.log('Connected successfully!')
    }

    socket.onclose = function () {
      console.log('Connection has been closed.')
    }

    socket.onmessage = function (e) {
      const msg = JSON.parse(e.data)

      switch (msg.type) {
        case 'get-canvas-state':
          if (!canvas.toDataURL()) return
          console.log('Sending canvas state')
          socket.send(
            JSON.stringify({
              type: 'canvas-state',
              data: canvas.toDataURL()
            })
          )
          break

        case 'canvas-state-from-server':
          console.log('Received the state')
          const img = new Image()
          img.src = msg.data
          img.onload = () => {
            console.log('image loaded')
            ctx.drawImage(img, 0, 0)
          }
          break

        case 'draw-line':
          console.log('draw-line', msg.data)
          if (!ctx) return console.log('No ctx here')
          drawLine({ ...msg?.data, ctx })
          break

        case 'clear':
          ctx.clearRect(0, 0, canvas.width, canvas.height)
          break

        default:
          console.log(msg.type)
      }
    }
  } catch (err) {
    console.error(err)
  }
}

initSocket()

canvas.addEventListener('mousedown', onMouseDown)
canvas.addEventListener('mousemove', onMouseMove)
canvas.addEventListener('mouseup', onMouseUp)

clearButton.addEventListener('click', () => {
  socket?.send(
    JSON.stringify({
      type: 'clear',
      data: null
    })
  )
})

/**
 *
 * @param {MouseEvent} e
 */
function onMouseDown (e) {
  drawing = true
  prevPoint = getMousePos(canvas, e)
}

/**
 *
 * @param {MouseEvent} e
 */
function onMouseMove (e) {
  if (!drawing) return

  const currentPoint = getMousePos(canvas, e)
  const color = colorPicker.value
  drawLine({ prevPoint, currentPoint, ctx, color })
  socket?.send(
    JSON.stringify({
      type: 'draw-line',
      data: { prevPoint, currentPoint, color }
    })
  )

  prevPoint = currentPoint
}

function onMouseUp () {
  drawing = false
}

/**
 *
 * @param {HTMLCanvasElement} canvas
 * @param {MouseEvent} e
 */
function getMousePos (canvas, e) {
  const rect = canvas.getBoundingClientRect()
  return {
    x: e.clientX - rect.left,
    y: e.clientY - rect.top
  }
}

function drawLine ({ prevPoint, currentPoint, ctx, color }) {
  console.log(prevPoint, currentPoint, color)
  ctx.beginPath()
  ctx.moveTo(prevPoint.x, prevPoint.y)
  ctx.lineTo(currentPoint.x, currentPoint.y)
  ctx.strokeStyle = color
  ctx.lineWidth = 2
  ctx.stroke()
  ctx.closePath()
}

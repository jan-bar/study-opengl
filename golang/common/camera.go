package common

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraMove int

const (
	ForWard  CameraMove = 1
	BackWard CameraMove = 2
	Left     CameraMove = 3
	Right    CameraMove = 4

	Yaw         float32 = -90.0
	Pitch       float32 = 0.0
	Speed       float32 = 2.5
	Sensitivity float32 = 0.1
	Zoom        float32 = 45.0
)

type Camera struct {
	// 相机属性
	Position mgl32.Vec3
	Front    mgl32.Vec3
	Up       mgl32.Vec3
	Right    mgl32.Vec3
	WorldUp  mgl32.Vec3
	// 欧拉角
	Yaw   float32
	Pitch float32
	// 相机选项
	MovementSpeed    float32
	MouseSensitivity float32
	Zoom             float32
}

type CameraOption func(*Camera)

func WithPosition(p mgl32.Vec3) CameraOption {
	return func(c *Camera) {
		c.Position = p
	}
}

func NewCamera(opts ...CameraOption) *Camera {
	camera := &Camera{
		Position:         mgl32.Vec3{0, 0, 0},
		Front:            mgl32.Vec3{0, 0, -1},
		WorldUp:          mgl32.Vec3{0, 1, 0},
		Yaw:              Yaw,
		Pitch:            Pitch,
		MovementSpeed:    Speed,
		MouseSensitivity: Sensitivity,
		Zoom:             Zoom,
	}

	for _, opt := range opts {
		opt(camera)
	}

	camera.updateCameraVectors()
	return camera
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

func (c *Camera) ProcessKeyboard(direction CameraMove, deltaTime float32) {
	// 处理从任何类似键盘的输入系统接收的输入。接受相机定义的 ENUM 形式的输入参数（将其从窗口系统中抽象出来）
	velocity := c.MovementSpeed * deltaTime
	switch direction {
	case ForWard:
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	case BackWard:
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	case Left:
		c.Position = c.Position.Sub(c.Right.Mul(velocity))
	case Right:
		c.Position = c.Position.Add(c.Right.Mul(velocity))
	default:
		panic("unexpected camera move")
	}
}

func (c *Camera) ProcessMouseMovement(xOffset, yOffset float32, constrainPitch ...bool) {
	// 处理从鼠标移动事件接收到的输入。只需要水平和垂直方向上的输入
	c.Yaw += xOffset * c.MouseSensitivity
	c.Pitch += yOffset * c.MouseSensitivity

	if len(constrainPitch) == 0 || constrainPitch[0] {
		// 确保当 pitch 超过 89.0 或 -89.0 度时，不会出现 flipped 相机
		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		}
		if c.Pitch < -89.0 {
			c.Pitch = -89.0
		}
	}

	// 使用更新的欧拉角更新前向量、右向量和上向量
	c.updateCameraVectors()
}

func (c *Camera) ProcessMouseScroll(yOffset float32) {
	// 处理从鼠标滚轮事件接收到的输入。只需要垂直轮轴上的输入
	c.Zoom -= yOffset
	if c.Zoom < 1.0 {
		c.Zoom = 1.0
	}
	if c.Zoom > 45.0 {
		c.Zoom = 45.0
	}
}

func (c *Camera) updateCameraVectors() {
	// 计算新的 Front 向量
	front := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(c.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch)))),
	}
	c.Front = front.Normalize()
	// 重新计算 Right 和 Up 向量
	// 对向量进行归一化，因为向上或向下看的次数越多，它们的长度就越接近 0，从而导致移动速度变慢。
	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}

package domain

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

// TransformTranslate
// 第一引数mは４行である必要がある
func TransformTranslate(m mat.Dense, x, y, z float64) mat.Dense {
	translateMatrix := mat.NewDense(4, 4, []float64{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	})

	var translated mat.Dense
	translated.Mul(translateMatrix, &m)
	return translated
}

// TransformRotate は3つの軸（X、Y、Z）での回転を適用します
// x, y, z はそれぞれの軸での回転角度（ラジアン）です
// 第一引数mは４行である必要がある
func TransformRotate(m mat.Dense, x, y, z float64) mat.Dense {
	// X軸回転行列
	rotateXMatrix := mat.NewDense(4, 4, []float64{
		1, 0, 0, 0,
		0, math.Cos(x), -math.Sin(x), 0,
		0, math.Sin(x), math.Cos(x), 0,
		0, 0, 0, 1,
	})

	// Y軸回転行列
	rotateYMatrix := mat.NewDense(4, 4, []float64{
		math.Cos(y), 0, math.Sin(y), 0,
		0, 1, 0, 0,
		-math.Sin(y), 0, math.Cos(y), 0,
		0, 0, 0, 1,
	})

	// Z軸回転行列
	rotateZMatrix := mat.NewDense(4, 4, []float64{
		math.Cos(z), -math.Sin(z), 0, 0,
		math.Sin(z), math.Cos(z), 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})

	// 回転を順番に適用: Z軸 -> Y軸 -> X軸の順で回転
	var temp1, temp2, rotated mat.Dense
	temp1.Mul(rotateZMatrix, &m)
	temp2.Mul(rotateYMatrix, &temp1)
	rotated.Mul(rotateXMatrix, &temp2)

	return rotated
}

// TransformParallelProjection 平行投影
// 第一引数mは４行である必要がある
func TransformParallelProjection(m mat.Dense) mat.Dense {
	projectionMatrix := mat.NewDense(4, 4, []float64{
		1, 0, 0, 0, // X軸
		0, 1, 0, 0, // Y軸
		0, 0, 0, 0, // Z軸（無視）
		0, 0, 0, 1, // 同次座標
	})

	var projected mat.Dense
	projected.Mul(projectionMatrix, &m)

	return projected
}

// transformScale は3次元オブジェクトの拡大・縮小変換を行います
// scaleX, scaleY, scaleZ はそれぞれのX、Y、Z軸方向の拡大率です
// 第一引数mは４行である必要がある
func TransformScale(m mat.Dense, scaleX, scaleY, scaleZ float64) mat.Dense {
	scaleMatrix := mat.NewDense(4, 4, []float64{
		scaleX, 0, 0, 0,
		0, scaleY, 0, 0,
		0, 0, scaleZ, 0,
		0, 0, 0, 1,
	})

	var scaled mat.Dense
	scaled.Mul(scaleMatrix, &m)
	return scaled
}

// transformScaleUniform は均等な拡大・縮小変換を行います
// scale は全軸方向の拡大率です
// 第一引数mは４行である必要がある
func TransformScaleUniform(m mat.Dense, scale float64) mat.Dense {
	return TransformScale(m, scale, scale, scale)
}

// TransformViewport はビューポート変換を行います
// SDL2の仕様に準拠した整数型の座標系に変換します
// SDL2の仕様：
// - 原点 (0,0) はウィンドウの左上
// - X軸は右方向が正
// - Y軸は下方向が正
//
// 第一引数mは４行である必要がある
// 実数値の単位座標(1.0f)を「画面の短辺」の{scaleRatio}%分、拡大する
func TransformViewport(m mat.Dense, width, height int32, scaleRatio float64) mat.Dense {
	// 短辺を基準にスケールを決める
	scale := math.Min(float64(width), float64(height)) * scaleRatio

	m = TransformScale(m, scale, -scale, 1) // SDL2に準拠するため上下反転する

	m = TransformTranslate(m, float64(width)/2, float64(height)/2, 0)

	return m
}

// T は行列を転置します
func T(m mat.Dense) mat.Dense {
	return *mat.DenseCopyOf((&m).T())
}

// CalcNormalFromPoints は3つの点から法線ベクトルを計算します
// 引数に渡す座標の順番で法線の表裏が変わるため注意する
// 法線方向は右ねじの法則で判別できます。
// 左手座標系なのでp2とp3を入れ替えてます。
func CalcNormalFromPoints(p1, p2, p3 Vector3D) Vector3D {
	p1v, p2v, p3v := p1.Vec(), p2.Vec(), p3.Vec()

	v1 := mat.NewVecDense(3, nil)
	v1.SubVec(&p3v, &p1v)

	v2 := mat.NewVecDense(3, nil)
	v2.SubVec(&p2v, &p1v)

	normal := CrossVecDense(*v1, *v2)

	normal = NormalizeVecDense(normal)

	return Vector3D(normal.RawVector().Data)
}

// CrossVecDense は2つのベクトルの外積を計算します
func CrossVecDense(a, b mat.VecDense) mat.VecDense {
	ax, ay, az := a.At(0, 0), a.At(1, 0), a.At(2, 0)
	bx, by, bz := b.At(0, 0), b.At(1, 0), b.At(2, 0)
	return *mat.NewVecDense(3, []float64{
		ay*bz - az*by,
		az*bx - ax*bz,
		ax*by - ay*bx,
	})
}

// NormalizeVecDense はベクトルを正規化します
func NormalizeVecDense(v mat.VecDense) mat.VecDense {
	norm := v.Norm(2)
	return *mat.NewVecDense(3, []float64{
		v.At(0, 0) / norm,
		v.At(1, 0) / norm,
		v.At(2, 0) / norm,
	})
}

// ClassifyEdgeByPlane は点が平面のどちら側にあるかを判定します
// 内側の場合はtrueを返します
// planeNormalは平面の法線ベクトル
// pInPlaneは平面の任意の点
func ClassifyEdgeByPlane(targetP Vector3D, planeNormal Vector3D, pInPlane Vector3D) bool {
	d := -(planeNormal[0]*pInPlane[0] + planeNormal[1]*pInPlane[1] + planeNormal[2]*pInPlane[2])
	result := planeNormal[0]*targetP[0] + planeNormal[1]*targetP[1] + planeNormal[2]*targetP[2] + d
	return result < 0
}

// IntersectPlaneIntersectionPoint は平面と線分の交点を計算します
func IntersectPlaneIntersectionPoint(planeNormal Vector3D, planePoint Vector3D, fromVertex, toVertex Vector3D) Vector3D {
	d := -(planeNormal[0]*planePoint[0] + planeNormal[1]*planePoint[1] + planeNormal[2]*planePoint[2])

	f := func(v Vector3D) float64 {
		return v[0]*planeNormal[0] + v[1]*planeNormal[1] + v[2]*planeNormal[2] + d
	}

	t := -f(fromVertex) / (f(toVertex) - f(fromVertex))

	p := fromVertex.Add(toVertex.Sub(fromVertex).Mul(t))

	return p
}

// Triangulate は多角形を三角形に分割します
func Triangulate(vertices []Vector3D) [][3]Vector3D {
	if len(vertices) < 3 {
		return [][3]Vector3D{}
	} else if len(vertices) == 3 {
		return [][3]Vector3D{{vertices[0], vertices[1], vertices[2]}}
	} else {
		triangles := make([][3]Vector3D, 0, len(vertices)-2)
		for i := 1; i < len(vertices)-1; i++ {
			triangles = append(triangles, [3]Vector3D{vertices[0], vertices[i], vertices[i+1]})
		}
		return triangles
	}
}

func CleanEdges(edges [][2]int) [][2]int {
	newEdges := make([][2]int, 0, len(edges))
	existMap := make(map[[2]int]bool, len(edges))

	makeKey := func(edge [2]int) [2]int {
		tmp := edge[:]
		sort.Ints(tmp)
		return [2]int{tmp[0], tmp[1]}
	}

	for _, edge := range edges {
		if edge[0] == edge[1] {
			// 同じ頂点への辺は破棄する
			continue
		}

		key := makeKey(edge)
		if _, ok := existMap[key]; ok {
			// 重複する辺は破棄する
			continue
		}

		existMap[key] = true
		newEdges = append(newEdges, edge)
	}

	return newEdges
}

func CleanTriangles(triangles [][3]int) [][3]int {
	newTriangles := make([][3]int, 0, len(triangles))
	existMap := make(map[[3]int]bool, len(triangles))

	makeKey := func(triangle [3]int) [3]int {
		tmp := triangle[:]
		sort.Ints(tmp)
		return [3]int{tmp[0], tmp[1], tmp[2]}
	}

	for _, triangle := range triangles {
		if triangle[0] == triangle[1] || triangle[1] == triangle[2] || triangle[2] == triangle[0] {
			// ３つの頂点の添字のうち、同じ添字を持っているものは破棄する
			continue
		}

		key := makeKey(triangle)
		if _, ok := existMap[key]; ok {
			// 重複する三角形は破棄する
			continue
		}

		existMap[key] = true
		newTriangles = append(newTriangles, triangle)
	}

	return newTriangles
}

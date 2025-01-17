package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestK8sQuantity(t *testing.T) {
	r1 := quantities{
		MilliCPU: resource.MustParse("1"),
		Memory:   resource.MustParse("100M"),
		ScalarResources: map[ResourceName]resource.Quantity{
			"nvidia.com/gpu": resource.MustParse("100Mi"),
		},
	}
	r2 := quantities{
		MilliCPU: resource.MustParse("2"),
		Memory:   resource.MustParse("2G"),
		ScalarResources: map[ResourceName]resource.Quantity{
			"nvidia.com/gpu": resource.MustParse("200M"),
		},
	}
	t.Logf(r1.Memory.String())
	mint, possible := r2.Memory.AsInt64()
	t.Logf("%d %v", mint, possible)

	// show unit
	r3 := quantities{
		MilliCPU: resource.MustParse("2"),
		Memory:   resource.MustParse("200000000000000"),
		ScalarResources: map[ResourceName]resource.Quantity{
			"nvidia.com/gpu": resource.MustParse("200Mi"),
		},
	}
	t.Logf(r3.Memory.String())

	// invalid value
	errMem, err := resource.ParseQuantity("xxx")
	if err != nil {
		t.Logf("parse quantity error: %v", err)
	}
	t.Logf(errMem.String())

	// add
	r1.MilliCPU.Add(r2.MilliCPU)
	r1.Memory.Add(r2.Memory)
	t.Logf("%s, %s", r1.MilliCPU.String(), r1.Memory.String())

	// sub
	r1.Memory.Sub(r2.Memory)
	t.Logf("%s, %s", r1.MilliCPU.String(), r1.Memory.String())

	r4, err := resource.ParseQuantity("")
	if err != nil {
		t.Logf("parse quantity error: %v", err)
	} else {
		t.Logf(r4.String())
	}

	rNegativem, err := resource.ParseQuantity("-1")
	if err != nil {
		t.Logf("parse quantity error: %v", err)
	} else {
		t.Logf(rNegativem.String())
	}
	_, err = resource.ParseQuantity("")
	assert.Error(t, err)
	t.Logf("%v", err)
}

func TestCheckResource(t *testing.T) {
	reses := []struct {
		res       string
		resType   string
		expectErr bool
	}{
		{"1", "cpu", false},
		{"-1", "cpu", true},
		{"0", "cpu", true},
		{"0.1", "cpu", false},
		{"", "cpu", true},

		{"1", "memory", false},
		{"-1", "memory", true},
		{"0", "memory", true},
		{"0.1Gi", "memory", false},

		{"1", "nvidia.com/gpu", false},
		{"0", "nvidia.com/gpu-xxx", false},
		{"-1", "xxx", true},
		{"0.1Gi", "nvidia.com/gpu", false},
	}

	for _, res := range reses {
		var err error
		switch res.resType {
		case "cpu":
			err = CheckCPUResource(res.res)
		case "memory":
			err = CheckMemoryResource(res.res)
		default:
			err = CheckScalarResource(res.res)
		}
		if res.expectErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestResourceInfo_LessEqual(t *testing.T) {
	queue1 := ResourceInfo{CPU: "1", Mem: "100M"}
	queue2 := ResourceInfo{
		CPU:             "1",
		Mem:             "100M",
		ScalarResources: ScalarResourcesType{"nvidia.com/gpu": "500"},
	}

	flavours := []ResourceInfo{
		{
			CPU:             "1",
			Mem:             "100M",
			ScalarResources: ScalarResourcesType{},
		},
		{
			CPU: "1",
			Mem: "100M",
			ScalarResources: ScalarResourcesType{
				"nvidia.com/gpu": "500",
			},
		},
	}

	cpuRes := flavours[0]

	gpuRes := flavours[1]

	// case1
	assert.True(t, cpuRes.LessEqual(queue1))
	assert.False(t, gpuRes.LessEqual(queue1))

	// case2
	assert.True(t, cpuRes.LessEqual(queue2))
	assert.True(t, gpuRes.LessEqual(queue2))
}

func TestResourceInfo_Sub(t *testing.T) {
	r1 := ResourceInfo{
		CPU:             "1",
		Mem:             "100M",
		ScalarResources: ScalarResourcesType{"nvidia.com/gpu": "500"},
	}
	r2 := ResourceInfo{
		CPU:             "1",
		Mem:             "100M",
		ScalarResources: ScalarResourcesType{"nvidia.com/gpu": "500"},
	}
	r3, err := r1.Sub(r2)
	if !assert.Nil(t, err) {
		t.Errorf(err.Error())
	}
	t.Logf("%v", r3)
}

func TestResourceInfo_Add(t *testing.T) {
	r1 := ResourceInfo{
		CPU:             "1",
		Mem:             "100M",
		ScalarResources: ScalarResourcesType{"nvidia.com/gpu": "500"},
	}
	r2 := ResourceInfo{
		CPU:             "1",
		Mem:             "100M",
		ScalarResources: ScalarResourcesType{"nvidia.com/gpu": "500"},
	}
	r3 := r1.Add(r2)

	t.Logf("%v", r3)
}

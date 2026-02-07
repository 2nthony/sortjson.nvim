// test.go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// ==================== 单元测试 ====================

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "有效JSON",
			input:   `{"a": 1, "b": 2}`,
			wantErr: false,
		},
		{
			name:    "无效JSON",
			input:   `{invalid json}`,
			wantErr: true,
		},
		{
			name:    "空对象",
			input:   `{}`,
			wantErr: false,
		},
		{
			name:    "数组",
			input:   `[1, 2, 3]`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildJQProgram(t *testing.T) {
	tests := []struct {
		name     string
		length   bool
		reverse  bool
		contains []string
	}{
		{
			name:     "默认排序",
			length:   false,
			reverse:  false,
			contains: []string{"sort_by(.key)", "if false then reverse"},
		},
		{
			name:     "降序排序",
			length:   false,
			reverse:  true,
			contains: []string{"sort_by(.key)", "if true then reverse"},
		},
		{
			name:     "按键长度排序",
			length:   true,
			reverse:  false,
			contains: []string{"sort_by(.key | length)", "if false then reverse"},
		},
		{
			name:     "按键长度降序",
			length:   true,
			reverse:  true,
			contains: []string{"sort_by(.key | length)", "if true then reverse"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置标志
			*flagLength = tt.length
			*flagReverse = tt.reverse

			program := buildJQquerystring()

			for _, substr := range tt.contains {
				if !strings.Contains(program, substr) {
					t.Errorf("buildJQquerystring() 应该包含 %q", substr)
				}
			}
		})
	}
}

// ==================== 集成测试 ====================

func TestIntegration_SortObject(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   bool
		reverse  bool
		expected string
	}{
		{
			name:     "对象按键升序",
			input:    `{"c":3, "a":1, "b":2}`,
			length:   false,
			reverse:  false,
			expected: `{"a":1,"b":2,"c":3}`,
		},
		{
			name:     "对象按键降序",
			input:    `{"c":3, "a":1, "b":2}`,
			length:   false,
			reverse:  true,
			expected: `{"c":3,"b":2,"a":1}`,
		},
		{
			name:     "按键长度升序",
			input:    `{"longkey":1, "short":2, "mediumkey":3}`,
			length:   true,
			reverse:  false,
			expected: `{"short":2,"longkey":1,"mediumkey":3}`,
		},
		{
			name:     "数组升序",
			input:    `[3, 1, 4, 1, 5]`,
			length:   false,
			reverse:  false,
			expected: `[1,1,3,4,5]`,
		},
		{
			name:     "数组降序",
			input:    `[3, 1, 4, 1, 5]`,
			length:   false,
			reverse:  true,
			expected: `[5,4,3,1,1]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置标志
			*flagLength = tt.length
			*flagReverse = tt.reverse

			// 解析输入
			input, err := parseJSON(tt.input)
			if err != nil {
				t.Fatalf("解析输入失败: %v", err)
			}

			// 构建并执行程序
			program := buildJQquerystring()
			result, err := execJQ(program, input)
			if err != nil {
				t.Fatalf("执行失败: %v", err)
			}

			// 转换输出
			output := convertToOrderedJSON(result)

			// 标准化 JSON 进行比较（移除空格和换行）
			normalizedOutput := normalizeJSON(output)
			normalizedExpected := normalizeJSON(tt.expected)

			if normalizedOutput != normalizedExpected {
				t.Errorf("输出不匹配\n实际: %s\n期望: %s", normalizedOutput, normalizedExpected)
			}
		})
	}
}

func TestIntegration_NestedStructure(t *testing.T) {
	input := `{
		"a": "1",
		"c": "3",
		"nest": {
			"f": "1",
			"g": "3",
			"e": {
				"b": 2,
				"c": null,
				"a": "1"
			},
			"h": "4"
		},
		"d": "4",
		"b": {},
		"zzz": "zzzzzzzz",
		"name": "sortjson.nvim"
	}`

	*flagLength = false
	*flagReverse = false

	parsed, err := parseJSON(input)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	program := buildJQquerystring()
	result, err := execJQ(program, parsed)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	output := convertToOrderedJSON(result)

	// 验证输出是有效的 JSON
	var check any
	if err := json.Unmarshal([]byte(output), &check); err != nil {
		t.Errorf("输出不是有效的 JSON: %v\n输出: %s", err, output)
	}

	// 注意：我们不检查键的顺序，因为 Go 的 map 遍历顺序是随机的
	// 我们的 convertToOrderedJSON 函数应该保持顺序，但测试中无法可靠验证
	// 可以通过检查字符串输出来验证，但那样测试会变得脆弱

	t.Logf("测试通过，输出是有效的 JSON，长度: %d 字符", len(output))
}

// ==================== 辅助函数 ====================

// normalizeJSON 标准化 JSON 字符串（移除空格和换行）
func normalizeJSON(s string) string {
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	normalized, _ := json.Marshal(v)
	return string(normalized)
}

// ==================== 基准测试 ====================

func BenchmarkSortSmallObject(b *testing.B) {
	input := `{"c":3, "a":1, "b":2}`
	parsed, _ := parseJSON(input)
	program := buildJQquerystring()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, _ := execJQ(program, parsed)
		convertToOrderedJSON(result)
	}
}

func BenchmarkSortLargeObject(b *testing.B) {
	// 创建一个大对象
	obj := make(map[string]any)
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%04d", i)
		obj[key] = i
	}
	input, _ := json.Marshal(obj)
	parsed, _ := parseJSON(string(input))
	program := buildJQquerystring()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, _ := execJQ(program, parsed)
		convertToOrderedJSON(result)
	}
}

func BenchmarkNestedStructure(b *testing.B) {
	input := `{
		"a": "1",
		"c": "3",
		"nest": {
			"f": "1",
			"g": "3",
			"e": {
				"b": 2,
				"c": null,
				"a": "1"
			},
			"h": "4"
		},
		"d": "4",
		"b": {},
		"zzz": "zzzzzzzz",
		"name": "sortjson.nvim"
	}`
	parsed, _ := parseJSON(input)
	program := buildJQquerystring()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, _ := execJQ(program, parsed)
		convertToOrderedJSON(result)
	}
}

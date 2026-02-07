package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/itchyny/gojq"
)

var (
	flagLength  = flag.Bool("length", false, "sort by key length")
	flagReverse = flag.Bool("reverse", false, "reverse sort (descending)")
)

func main() {
	flag.Parse()

	args := flag.Args()

	jsonStr := strings.Join(args, " ")

	input, err := parseJSON(jsonStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	qs := buildJQquerystring()

	query, err := execJQ(qs, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	output := convertToOrderedJSON(query)
	fmt.Println(output)
}

func parseJSON(jsonStr string) (any, error) {
	var input any
	if err := json.Unmarshal([]byte(jsonStr), &input); err != nil {
		return nil, fmt.Errorf("Parse failed: %v\ninput: %s", err, jsonStr)
	}
	return input, nil
}

func buildJQquerystring() string {
	var sortExpr = "sort_by(.key)"
	if *flagLength {
		sortExpr = "sort_by(.key | length)"
	}

	return fmt.Sprintf(`
def sort_recursive:
  if type == "object" then
    # 递归处理所有值
    with_entries(.value |= sort_recursive)
    # 按键排序并转换为条目数组
    | to_entries | %s
    | if %v then reverse else . end
  elif type == "array" then
    # 对数组排序
    sort
    | if %v then reverse else . end
    | map(sort_recursive)
  else
    .
  end;

sort_recursive
`, sortExpr, *flagReverse, *flagReverse)
}

func execJQ(qs string, input any) (any, error) {
	query, err := gojq.Parse(qs)
	if err != nil {
		return nil, fmt.Errorf("gojq program %v\n", err)
	}

	iter := query.Run(input)
	v, ok := iter.Next()
	if !ok {
		return nil, fmt.Errorf("no output")
	}
	if err, isErr := v.(error); isErr {
		return nil, fmt.Errorf("Failed: %v\n", err)
	}

	return v, nil
}

// 递归转换任何值为有序 JSON
func convertToOrderedJSON(v any) string {
	return convertToOrderedJSONIndent(v, "", "  ")
}

func convertToOrderedJSONIndent(v any, indent, indentStep string) string {
	switch val := v.(type) {
	case []any:
		// 检查是否是条目数组
		if isEntriesArray(val) {
			// 转换为有序对象
			return entriesToOrderedObject(val, indent, indentStep)
		}
		// 普通数组
		return arrayToJSON(val, indent, indentStep)
	case map[string]any:
		// 普通对象（不应该出现，因为 jq 输出的是条目数组）
		return objectToJSON(val, indent, indentStep)
	default:
		// 基本类型
		jsonBytes, _ := json.Marshal(v)
		return string(jsonBytes)
	}
}

// 检查是否是条目数组
func isEntriesArray(arr []any) bool {
	if len(arr) == 0 {
		return false
	}
	for _, item := range arr {
		if m, ok := item.(map[string]any); ok {
			if _, hasKey := m["key"]; !hasKey {
				return false
			}
			if _, hasValue := m["value"]; !hasValue {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// 将条目数组转换为有序对象
func entriesToOrderedObject(entries []any, indent, indentStep string) string {
	if len(entries) == 0 {
		return "{}"
	}

	var buf strings.Builder
	buf.WriteString("{")

	nextIndent := indent + indentStep

	for i, entry := range entries {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
		buf.WriteString(nextIndent)

		if m, ok := entry.(map[string]any); ok {
			// 编码 key
			if key, ok := m["key"].(string); ok {
				keyJSON, _ := json.Marshal(key)
				buf.WriteString(string(keyJSON))
				buf.WriteString(": ")

				// 递归编码 value
				value := m["value"]
				valueStr := convertToOrderedJSONIndent(value, nextIndent, indentStep)
				buf.WriteString(valueStr)
			}
		}
	}

	buf.WriteString("\n")
	buf.WriteString(indent)
	buf.WriteString("}")
	return buf.String()
}

// 普通数组转 JSON
func arrayToJSON(arr []any, indent, indentStep string) string {
	if len(arr) == 0 {
		return "[]"
	}

	var buf strings.Builder
	buf.WriteString("[")

	nextIndent := indent + indentStep

	for i, item := range arr {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
		buf.WriteString(nextIndent)

		itemStr := convertToOrderedJSONIndent(item, nextIndent, indentStep)
		buf.WriteString(itemStr)
	}

	buf.WriteString("\n")
	buf.WriteString(indent)
	buf.WriteString("]")
	return buf.String()
}

// 普通对象转 JSON（回退方案）
func objectToJSON(obj map[string]any, indent, indentStep string) string {
	if len(obj) == 0 {
		return "{}"
	}

	// 获取排序后的键
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}

	var buf strings.Builder
	buf.WriteString("{")

	nextIndent := indent + indentStep

	for i, key := range keys {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
		buf.WriteString(nextIndent)

		keyJSON, _ := json.Marshal(key)
		buf.WriteString(string(keyJSON))
		buf.WriteString(": ")

		value := obj[key]
		valueStr := convertToOrderedJSONIndent(value, nextIndent, indentStep)
		buf.WriteString(valueStr)
	}

	buf.WriteString("\n")
	buf.WriteString(indent)
	buf.WriteString("}")
	return buf.String()
}

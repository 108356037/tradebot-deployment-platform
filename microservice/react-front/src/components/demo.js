import { useState, useRef, useCallback } from "react";
import Cron from "react-js-cron";
import { Input, Divider } from "antd";
import 'antd/dist/antd.css';

function Demo() {
  const inputRef = useRef(null)
  const defaultValue = '30 5 * * 1,6'
  const [cronValue, setCronValue] = useState(defaultValue)
  const customSetValue = useCallback(
    (newValue) => {
      setCronValue(newValue)
      inputRef.current?.setValue(newValue)
    },
    [inputRef]
  )
  const [error, onError] = useState({"type": "invalid_cron", "description": ""})

  return (
    <div>
      <Input
        ref={inputRef}
        onBlur={(event) => {
          setCronValue(event.target.value)
        }}
        onPressEnter={() => {
          setCronValue(inputRef.current?.input.value || '')
        }}
      />

      <Divider>OR</Divider>

      <Cron className="react-js-cron" value={cronValue} setValue={customSetValue} onError={onError}/>

      <p style={{ marginTop: 20 }}>
        Error: {error ? error.description : 'undefined'}
      </p>
    </div>
  )
}

export default Demo;

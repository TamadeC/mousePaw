import { useState, useEffect, useRef } from 'react';
import { GetConfig, UpdateConfig, GetStatus, Start, Stop, Pause, Resume, GetLogs, StartRecording, StopRecording, GetRecordingStatus, GetRecording, ClearRecording, SaveRecording, LoadRecording, ListRecordings, DeleteRecording } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime/runtime';
import { config as WailsModels } from '../wailsjs/go/models';
import { recorder as RecorderModels } from '../wailsjs/go/models';

interface Config {
  operation_type: string;
  move_interval: number;
  move_random: boolean;
  click_interval: number;
  click_type: string;
  click_count: number;
  scroll_interval: number;
  scroll_dir: string;
  scroll_amount: number;
  type_interval: number;
  type_text: string;
  replay_interval: number;
  replay_repeat: boolean;
  replay_file: string;
  auto_start: boolean;
  minimize_to_tray: boolean;
  hotkeys: {
    start: string;
    stop: string;
    pause: string;
  };
}

interface LogEntry {
  time: string;
  level: string;
  message: string;
}

type RecorderStatusInfo = RecorderModels.RecorderStatusInfo;
type RecordedAction = RecorderModels.RecordedAction;
type Recording = RecorderModels.Recording;

type Status = 'stopped' | 'running' | 'paused';

const defaultConfig: Config = {
  operation_type: 'move',
  move_interval: 5,
  move_random: true,
  click_interval: 3,
  click_type: 'left',
  click_count: 1,
  scroll_interval: 5,
  scroll_dir: 'down',
  scroll_amount: 3,
  type_interval: 1,
  type_text: '',
  replay_interval: 30,
  replay_repeat: false,
  replay_file: '',
  auto_start: false,
  minimize_to_tray: true,
  hotkeys: {
    start: 'ctrl+f6',
    stop: 'ctrl+f7',
    pause: 'ctrl+f8',
  },
};

export default function App() {
  const [config, setConfig] = useState<Config>(defaultConfig);
  const [status, setStatus] = useState<Status>('stopped');
  const [activeTab, setActiveTab] = useState<'settings' | 'system' | 'logs' | 'recording'>('settings');
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const logsEndRef = useRef<HTMLDivElement>(null);

  const [recorderStatus, setRecorderStatus] = useState<RecorderStatusInfo>({ status: 'idle', count: 0, duration: 0 });
  const [recording, setRecording] = useState<Recording | null>(null);
  const [savedFiles, setSavedFiles] = useState<string[]>([]);
  const [saveName, setSaveName] = useState('');

  useEffect(() => {
    loadConfig();
    loadStatus();
    loadLogs();
    loadRecordingStatus();
    loadSavedFiles();

    const cleanupStatus = EventsOn('statusChanged', (s: string) => {
      setStatus(s as Status);
    });

    const cleanupLog = EventsOn('newLog', (entry: LogEntry) => {
      setLogs(prev => [...prev, entry]);
    });

    const cleanupRecorder = EventsOn('recorderStatus', (info: RecorderStatusInfo) => {
      setRecorderStatus(info);
    });

    return () => {
      cleanupStatus();
      cleanupLog();
      cleanupRecorder();
    };
  }, []);

  useEffect(() => {
    if (activeTab === 'logs') {
      logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }
  }, [logs, activeTab]);

  async function loadConfig() {
    try {
      const cfg = await GetConfig();
      setConfig(cfg);
    } catch (e) {
      console.error('Failed to load config:', e);
    }
  }

  async function loadStatus() {
    try {
      const s = await GetStatus();
      setStatus(s as Status);
    } catch (e) {
      console.error('Failed to load status:', e);
    }
  }

  async function loadLogs() {
    try {
      const l = await GetLogs();
      setLogs(l);
    } catch (e) {
      console.error('Failed to load logs:', e);
    }
  }

  async function loadRecordingStatus() {
    try {
      const s = await GetRecordingStatus();
      setRecorderStatus(s);
    } catch (e) {
      console.error('Failed to load recorder status:', e);
    }
  }

  async function loadSavedFiles() {
    try {
      const files = await ListRecordings();
      setSavedFiles(files);
    } catch (e) {
      console.error('Failed to load recordings:', e);
    }
  }

  async function updateConfig(key: keyof Config, value: any) {
    const newConfig = { ...config, [key]: value };
    setConfig(newConfig);
    try {
      await UpdateConfig(WailsModels.Config.createFrom(newConfig));
    } catch (e) {
      console.error('Failed to save config:', e);
    }
  }

  async function handleStart() {
    try {
      await Start();
    } catch (e) {
      console.error('Failed to start:', e);
    }
  }

  async function handleStop() {
    try {
      await Stop();
    } catch (e) {
      console.error('Failed to stop:', e);
    }
  }

  async function handlePause() {
    try {
      await Pause();
    } catch (e) {
      console.error('Failed to pause:', e);
    }
  }

  async function handleResume() {
    try {
      await Resume();
    } catch (e) {
      console.error('Failed to resume:', e);
    }
  }

  async function handleStartRecording() {
    try {
      await StartRecording();
    } catch (e) {
      console.error('Failed to start recording:', e);
    }
  }

  async function handleStopRecording() {
    try {
      const rec = await StopRecording();
      if (rec) setRecording(rec);
      await loadSavedFiles();
    } catch (e) {
      console.error('Failed to stop recording:', e);
    }
  }

  async function handleClearRecording() {
    try {
      await ClearRecording();
      setRecording(null);
    } catch (e) {
      console.error('Failed to clear recording:', e);
    }
  }

  async function handleSaveRecording() {
    if (!saveName.trim()) return;
    try {
      await SaveRecording(saveName.trim());
      setSaveName('');
      await loadSavedFiles();
    } catch (e) {
      console.error('Failed to save recording:', e);
    }
  }

  async function handleLoadRecording(name: string) {
    try {
      const rec = await LoadRecording(name);
      if (rec) setRecording(rec);
    } catch (e) {
      console.error('Failed to load recording:', e);
    }
  }

  async function handleDeleteRecording(name: string) {
    try {
      await DeleteRecording(name);
      await loadSavedFiles();
    } catch (e) {
      console.error('Failed to delete recording:', e);
    }
  }

  const statusColors = {
    stopped: 'bg-gray-500',
    running: 'bg-green-500',
    paused: 'bg-yellow-500',
  };

  const statusLabels = {
    stopped: '已停止',
    running: '运行中',
    paused: '已暂停',
  };

  return (
    <div className="h-screen flex flex-col bg-[#0f172a] text-gray-100 select-none">
      {/* Title Bar */}
      <div className="flex items-center justify-between px-4 py-3 bg-[#1e293b] border-b border-[#334155]">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-indigo-500 rounded-lg flex items-center justify-center">
            <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
            </svg>
          </div>
          <h1 className="text-lg font-semibold">MousePaw</h1>
        </div>
        <div className="flex items-center gap-2">
          <span className={`w-2.5 h-2.5 rounded-full ${statusColors[status]} animate-pulse`} />
          <span className="text-sm text-gray-400">{statusLabels[status]}</span>
        </div>
        <div className="flex items-center gap-2">
          {status === 'stopped' && (
            <button
              onClick={handleStart}
              className="px-3 py-1.5 bg-green-500 hover:bg-green-600 text-white text-sm rounded-lg transition-colors"
            >
              开始
            </button>
          )}
          {status === 'running' && (
            <>
              <button
                onClick={handlePause}
                className="px-3 py-1.5 bg-yellow-500 hover:bg-yellow-600 text-white text-sm rounded-lg transition-colors"
              >
                暂停
              </button>
              <button
                onClick={handleStop}
                className="px-3 py-1.5 bg-red-500 hover:bg-red-600 text-white text-sm rounded-lg transition-colors"
              >
                停止
              </button>
            </>
          )}
          {status === 'paused' && (
            <>
              <button
                onClick={handleResume}
                className="px-3 py-1.5 bg-green-500 hover:bg-green-600 text-white text-sm rounded-lg transition-colors"
              >
                恢复
              </button>
              <button
                onClick={handleStop}
                className="px-3 py-1.5 bg-red-500 hover:bg-red-600 text-white text-sm rounded-lg transition-colors"
              >
                停止
              </button>
            </>
          )}
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-hidden flex flex-col">
        {/* Tabs */}
        <div className="flex border-b border-[#334155]">
          {[
            { key: 'settings', label: '操作设置' },
            { key: 'recording', label: '操作录制' },
            { key: 'system', label: '系统设置' },
            { key: 'logs', label: '执行日志' },
          ].map(tab => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key as any)}
              className={`flex-1 py-3 text-sm font-medium transition-all ${
                activeTab === tab.key
                  ? 'text-indigo-400 border-b-2 border-indigo-400 bg-indigo-500/10'
                  : 'text-gray-400 hover:text-gray-200 hover:bg-[#1e293b]/50'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Tab Content */}
        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          {activeTab === 'settings' && (
            <OperationSettings config={config} updateConfig={updateConfig} recording={recording} />
          )}
          {activeTab === 'recording' && (
            <RecordingPanel
              recorderStatus={recorderStatus}
              recording={recording}
              savedFiles={savedFiles}
              saveName={saveName}
              setSaveName={setSaveName}
              onStartRecording={handleStartRecording}
              onStopRecording={handleStopRecording}
              onClearRecording={handleClearRecording}
              onSaveRecording={handleSaveRecording}
              onLoadRecording={handleLoadRecording}
              onDeleteRecording={handleDeleteRecording}
            />
          )}
          {activeTab === 'system' && (
            <SystemSettings config={config} updateConfig={updateConfig} />
          )}
          {activeTab === 'logs' && (
            <LogPanel logs={logs} logsEndRef={logsEndRef} />
          )}
        </div>

        {/* Hotkey Hint */}
        <div className="px-4 py-3 bg-[#1e293b]/50 border-t border-[#334155]">
          <div className="flex items-center justify-center gap-6 text-sm">
            <div className="flex items-center gap-2">
              <kbd className="px-2 py-1 bg-[#334155] rounded text-indigo-400 font-mono text-xs">
                {config.hotkeys.start.split('+').join(' + ')}
              </kbd>
              <span className="text-gray-400">开始</span>
            </div>
            <div className="flex items-center gap-2">
              <kbd className="px-2 py-1 bg-[#334155] rounded text-indigo-400 font-mono text-xs">
                {config.hotkeys.pause.split('+').join(' + ')}
              </kbd>
              <span className="text-gray-400">暂停/恢复</span>
            </div>
            <div className="flex items-center gap-2">
              <kbd className="px-2 py-1 bg-[#334155] rounded text-indigo-400 font-mono text-xs">
                {config.hotkeys.stop.split('+').join(' + ')}
              </kbd>
              <span className="text-gray-400">停止</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function Card({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="bg-[#1e293b] rounded-xl border border-[#334155] overflow-hidden">
      <div className="px-4 py-3 border-b border-[#334155]">
        <h3 className="text-sm font-medium text-gray-300">{title}</h3>
      </div>
      <div className="p-4 space-y-4">
        {children}
      </div>
    </div>
  );
}

function Toggle({ label, checked, onChange }: { label: string; checked: boolean; onChange: (v: boolean) => void }) {
  return (
    <label className="flex items-center justify-between cursor-pointer">
      <span className="text-sm text-gray-300">{label}</span>
      <div
        className={`w-11 h-6 rounded-full transition-colors ${checked ? 'bg-indigo-500' : 'bg-[#475569]'}`}
        onClick={() => onChange(!checked)}
      >
        <div
          className={`w-5 h-5 rounded-full bg-white shadow-md transform transition-transform mt-0.5 ${checked ? 'translate-x-[22px]' : 'translate-x-0.5'}`}
        />
      </div>
    </label>
  );
}

function Slider({ label, value, min, max, step, unit, onChange }: {
  label: string; value: number; min: number; max: number; step: number; unit: string;
  onChange: (v: number) => void;
}) {
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = parseFloat(e.target.value);
    if (!isNaN(val) && val >= min && val <= max) {
      onChange(val);
    }
  };

  return (
    <div className="space-y-2">
      <div className="flex justify-between items-center">
        <span className="text-sm text-gray-300">{label}</span>
        <div className="flex items-center gap-2">
          <input
            type="number"
            min={min}
            max={max}
            step={step}
            value={value}
            onChange={handleInputChange}
            className="w-16 bg-[#334155] text-indigo-400 text-sm rounded-lg px-2 py-1 border border-[#475569] focus:border-indigo-500 focus:outline-none text-center font-medium"
          />
          <span className="text-sm text-gray-400">{unit}</span>
        </div>
      </div>
      <input
        type="range"
        min={min}
        max={max}
        step={step}
        value={value}
        onChange={e => onChange(Number(e.target.value))}
        className="w-full h-2 bg-[#334155] rounded-lg appearance-none cursor-pointer accent-indigo-500"
      />
    </div>
  );
}

function Select({ label, value, options, onChange }: {
  label: string; value: string; options: { value: string; label: string }[];
  onChange: (v: string) => void;
}) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-sm text-gray-300">{label}</span>
      <select
        value={value}
        onChange={e => onChange(e.target.value)}
        className="bg-[#334155] text-gray-200 text-sm rounded-lg px-3 py-1.5 border border-[#475569] focus:border-indigo-500 focus:outline-none"
      >
        {options.map(opt => (
          <option key={opt.value} value={opt.value}>{opt.label}</option>
        ))}
      </select>
    </div>
  );
}

function NumberInput({ label, value, min, max, onChange }: {
  label: string; value: number; min: number; max: number;
  onChange: (v: number) => void;
}) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-sm text-gray-300">{label}</span>
      <input
        type="number"
        min={min}
        max={max}
        value={value}
        onChange={e => onChange(Number(e.target.value))}
        className="w-20 bg-[#334155] text-gray-200 text-sm rounded-lg px-3 py-1.5 border border-[#475569] focus:border-indigo-500 focus:outline-none text-center"
      />
    </div>
  );
}

function OperationTypeSelector({ value, onChange }: { value: string; onChange: (v: string) => void }) {
  const options = [
    { key: 'move', label: '鼠标移动', icon: 'M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122' },
    { key: 'click', label: '鼠标点击', icon: 'M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5' },
    { key: 'scroll', label: '滚轮滚动', icon: 'M19 14l-7 7m0 0l-7-7m7 7V3' },
    { key: 'type', label: '键盘输入', icon: 'M12 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2M3 12l6.414 6.414a2 2 0 001.414.586H19a2 2 0 002-2V7a2 2 0 00-2-2h-8.172a2 2 0 00-1.414.586L3 12z' },
    { key: 'replay', label: '操作回放', icon: 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15' },
  ];

  return (
    <Card title="操作类型">
      <div className="grid grid-cols-5 gap-2">
        {options.map(opt => (
          <button
            key={opt.key}
            onClick={() => onChange(opt.key)}
            className={`flex flex-col items-center gap-2 p-4 rounded-xl border-2 transition-all ${
              value === opt.key
                ? 'border-indigo-500 bg-indigo-500/20 text-indigo-400'
                : 'border-[#475569] bg-[#1e293b] text-gray-400 hover:border-gray-500 hover:text-gray-300'
            }`}
          >
            <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={opt.icon} />
            </svg>
            <span className="text-sm font-medium">{opt.label}</span>
          </button>
        ))}
      </div>
    </Card>
  );
}

function OperationSettings({ config, updateConfig, recording }: { config: Config; updateConfig: (k: keyof Config, v: any) => void; recording: Recording | null }) {
  return (
    <>
      <OperationTypeSelector
        value={config.operation_type}
        onChange={v => updateConfig('operation_type', v)}
      />
      {config.operation_type === 'move' && (
        <Card title="移动参数">
          <Slider
            label="移动间隔"
            value={config.move_interval}
            min={1}
            max={60}
            step={0.5}
            unit="秒"
            onChange={v => updateConfig('move_interval', v)}
          />
          <Toggle
            label="随机位置移动"
            checked={config.move_random}
            onChange={v => updateConfig('move_random', v)}
          />
          <p className="text-xs text-gray-500">
            {config.move_random
              ? '鼠标将移动到屏幕上的随机位置'
              : '鼠标将向随机方向移动100-200像素'}
          </p>
        </Card>
      )}
      {config.operation_type === 'click' && (
        <Card title="点击参数">
          <Slider
            label="点击间隔"
            value={config.click_interval}
            min={0.5}
            max={30}
            step={0.5}
            unit="秒"
            onChange={v => updateConfig('click_interval', v)}
          />
          <Select
            label="点击按钮"
            value={config.click_type}
            options={[
              { value: 'left', label: '左键' },
              { value: 'right', label: '右键' },
              { value: 'middle', label: '中键' },
            ]}
            onChange={v => updateConfig('click_type', v)}
          />
          <NumberInput
            label="连续点击次数"
            value={config.click_count}
            min={1}
            max={10}
            onChange={v => updateConfig('click_count', v)}
          />
        </Card>
      )}
      {config.operation_type === 'scroll' && (
        <Card title="滚动参数">
          <Slider
            label="滚动间隔"
            value={config.scroll_interval}
            min={1}
            max={30}
            step={0.5}
            unit="秒"
            onChange={v => updateConfig('scroll_interval', v)}
          />
          <Select
            label="滚动方向"
            value={config.scroll_dir}
            options={[
              { value: 'up', label: '向上' },
              { value: 'down', label: '向下' },
              { value: 'left', label: '向左' },
              { value: 'right', label: '向右' },
            ]}
            onChange={v => updateConfig('scroll_dir', v)}
          />
          <Slider
            label="滚动幅度"
            value={config.scroll_amount}
            min={1}
            max={20}
            step={1}
            unit="格"
            onChange={v => updateConfig('scroll_amount', v)}
          />
        </Card>
      )}
      {config.operation_type === 'type' && (
        <Card title="键盘输入参数">
          <Slider
            label="输入间隔"
            value={config.type_interval}
            min={0.5}
            max={10}
            step={0.5}
            unit="秒"
            onChange={v => updateConfig('type_interval', v)}
          />
          <div className="space-y-2">
            <label className="text-sm text-gray-300">输入文本</label>
            <textarea
              value={config.type_text}
              onChange={e => updateConfig('type_text', e.target.value)}
              placeholder="输入要自动键入的文本..."
              className="w-full bg-[#334155] text-gray-200 text-sm rounded-lg px-3 py-2 border border-[#475569] focus:border-indigo-500 focus:outline-none h-24 resize-none"
            />
          </div>
          <p className="text-xs text-gray-500">
            将按照设定的间隔重复输入上方文本
          </p>
        </Card>
      )}
      {config.operation_type === 'replay' && (
        <Card title="回放参数">
          <Slider
            label="回放间隔"
            value={config.replay_interval}
            min={1}
            max={300}
            step={1}
            unit="秒"
            onChange={v => updateConfig('replay_interval', v)}
          />
          <Toggle
            label="循环回放"
            checked={config.replay_repeat}
            onChange={v => updateConfig('replay_repeat', v)}
          />
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-300">当前录制文件</span>
            <span className="text-sm text-indigo-400 font-mono">
              {config.replay_file || '未选择'}
            </span>
          </div>
          {recording && recording.actions.length > 0 && (
            <button
              onClick={() => {
                updateConfig('replay_file', config.replay_file || 'recording');
              }}
              className="w-full py-2 bg-indigo-500/20 text-indigo-400 text-sm rounded-lg border border-indigo-500/30 hover:bg-indigo-500/30 transition-colors"
            >
              使用当前录制 ({recording.actions.length} 个动作)
            </button>
          )}
          <p className="text-xs text-gray-500">
            {config.replay_repeat
              ? '录制操作将循环回放，直到手动停止'
              : `每隔 ${config.replay_interval} 秒触发一次完整回放`}
          </p>
        </Card>
      )}
    </>
  );
}

function SystemSettings({ config, updateConfig }: { config: Config; updateConfig: (k: keyof Config, v: any) => void }) {
  const updateHotkeys = (key: string, value: string) => {
    const newHotkeys = { ...config.hotkeys, [key]: value };
    updateConfig('hotkeys', newHotkeys);
  };

  return (
    <>
      <Card title="系统选项">
        <Toggle
          label="开机自启动"
          checked={config.auto_start}
          onChange={v => updateConfig('auto_start', v)}
        />
        <Toggle
          label="关闭时最小化到托盘"
          checked={config.minimize_to_tray}
          onChange={v => updateConfig('minimize_to_tray', v)}
        />
      </Card>
      <Card title="快捷键设置">
        <HotkeyInput
          label="开始"
          value={config.hotkeys.start}
          onChange={v => updateHotkeys('start', v)}
        />
        <HotkeyInput
          label="停止"
          value={config.hotkeys.stop}
          onChange={v => updateHotkeys('stop', v)}
        />
        <HotkeyInput
          label="暂停/恢复"
          value={config.hotkeys.pause}
          onChange={v => updateHotkeys('pause', v)}
        />
        <p className="text-xs text-gray-500">
          点击输入框后按下新的快捷键组合，修改后需重启应用生效
        </p>
      </Card>
      <Card title="关于">
        <div className="text-center space-y-2 py-2">
          <div className="w-12 h-12 bg-indigo-500 rounded-xl mx-auto flex items-center justify-center">
            <svg className="w-8 h-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
            </svg>
          </div>
          <p className="text-sm font-medium text-gray-200">MousePaw v1.1.0</p>
          <p className="text-xs text-gray-400">鼠标自动化工具</p>
          <p className="text-xs text-gray-500">按 {config.hotkeys.start} 开始，按 {config.hotkeys.stop} 停止</p>
        </div>
      </Card>
    </>
  );
}

function HotkeyInput({ label, value, onChange }: { label: string; value: string; onChange: (v: string) => void }) {
  const [recording, setRecording] = useState(false);
  const [keys, setKeys] = useState<string[]>([]);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    e.preventDefault();
    if (!recording) return;

    const key = e.key.toLowerCase();
    if (key === 'control' || key === 'shift' || key === 'alt' || key === 'meta') {
      return;
    }

    const modifiers: string[] = [];
    if (e.ctrlKey) modifiers.push('ctrl');
    if (e.shiftKey) modifiers.push('shift');
    if (e.altKey) modifiers.push('alt');

    const hotkey = [...modifiers, key].join('+');
    onChange(hotkey);
    setRecording(false);
    setKeys([]);
  };

  const startRecording = () => {
    setRecording(true);
    setKeys([]);
  };

  const displayValue = value.split('+').map(k => k.charAt(0).toUpperCase() + k.slice(1)).join(' + ');

  return (
    <div className="flex items-center justify-between">
      <span className="text-sm text-gray-300">{label}</span>
      <button
        onClick={startRecording}
        onKeyDown={handleKeyDown}
        className={`px-3 py-1.5 rounded-lg text-sm font-mono transition-all ${
          recording
            ? 'bg-indigo-500/20 border-indigo-500 text-indigo-400'
            : 'bg-[#334155] border-[#475569] text-gray-200 hover:border-indigo-500'
        } border`}
      >
        {recording ? '按下快捷键...' : displayValue}
      </button>
    </div>
  );
}

function LogPanel({ logs, logsEndRef }: { logs: LogEntry[]; logsEndRef: React.RefObject<HTMLDivElement> }) {
  const levelColors: Record<string, string> = {
    INFO: 'text-blue-400',
    ERROR: 'text-red-400',
  };

  return (
    <Card title="执行日志">
      <div className="bg-[#0f172a] rounded-lg p-3 h-[calc(100vh-300px)] overflow-y-auto font-mono text-xs">
        {logs.length === 0 ? (
          <div className="text-gray-500 text-center py-8">暂无日志</div>
        ) : (
          logs.map((log, i) => (
            <div key={i} className="py-1 flex gap-2">
              <span className="text-gray-500 shrink-0">{log.time}</span>
              <span className={`shrink-0 ${levelColors[log.level] || 'text-gray-400'}`}>
                [{log.level}]
              </span>
              <span className="text-gray-300">{log.message}</span>
            </div>
          ))
        )}
        <div ref={logsEndRef} />
      </div>
    </Card>
  );
}

function RecordingPanel({
  recorderStatus, recording, savedFiles, saveName, setSaveName,
  onStartRecording, onStopRecording, onClearRecording,
  onSaveRecording, onLoadRecording, onDeleteRecording,
}: {
  recorderStatus: RecorderStatusInfo;
  recording: Recording | null;
  savedFiles: string[];
  saveName: string;
  setSaveName: (v: string) => void;
  onStartRecording: () => void;
  onStopRecording: () => void;
  onClearRecording: () => void;
  onSaveRecording: () => void;
  onLoadRecording: (name: string) => void;
  onDeleteRecording: (name: string) => void;
}) {
  const isRecording = recorderStatus.status === 'recording';

  const actionTypeLabel = (t: string) => {
    switch (t) {
      case 'move': return '移动';
      case 'click': return '点击';
      case 'scroll': return '滚轮';
      case 'key': return '按键';
      default: return t;
    }
  };

  const actionDetail = (a: RecordedAction) => {
    switch (a.type) {
      case 'move': return `(${a.x ?? 0}, ${a.y ?? 0})`;
      case 'click': {
        const btn = { left: '左键', right: '右键', center: '中键' }[a.button || ''] || a.button;
        return `${btn} (${a.x ?? 0}, ${a.y ?? 0})`;
      }
      case 'scroll': {
        const dir = { up: '↑', down: '↓', left: '←', right: '→' }[a.direction || ''] || a.direction;
        return `${dir} ${a.amount ?? 0}格`;
      }
      case 'key': return a.keychar || '';
      default: return '';
    }
  };

  const typeColor = (t: string) => {
    switch (t) {
      case 'move': return 'text-blue-400';
      case 'click': return 'text-green-400';
      case 'scroll': return 'text-yellow-400';
      case 'key': return 'text-purple-400';
      default: return 'text-gray-400';
    }
  };

  return (
    <>
      {/* Recording Controls */}
      <Card title="录制控制">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <span className={`w-3 h-3 rounded-full ${isRecording ? 'bg-red-500 animate-pulse' : 'bg-gray-500'}`} />
            <span className={`text-sm font-medium ${isRecording ? 'text-red-400' : 'text-gray-400'}`}>
              {isRecording ? '录制中...' : '就绪'}
            </span>
          </div>
          <div className="flex items-center gap-2">
            {isRecording ? (
              <button
                onClick={onStopRecording}
                className="px-4 py-2 bg-red-500 hover:bg-red-600 text-white text-sm rounded-lg transition-colors"
              >
                停止录制
              </button>
            ) : (
              <button
                onClick={onStartRecording}
                className="px-4 py-2 bg-indigo-500 hover:bg-indigo-600 text-white text-sm rounded-lg transition-colors"
              >
                开始录制
              </button>
            )}
          </div>
        </div>
        <div className="flex gap-4 text-sm text-gray-400 pt-2">
          <span>动作数: <span className="text-gray-200 font-mono">{recorderStatus.count}</span></span>
          <span>时长: <span className="text-gray-200 font-mono">{(recorderStatus.duration ?? 0).toFixed(1)}s</span></span>
        </div>
        <p className="text-xs text-gray-500">
          {isRecording
            ? '正在录制全局鼠标和键盘操作，所有操作将被记录供回放使用'
            : '点击"开始录制"后，操作鼠标（移动、点击、滚轮）和键盘，程序将自动记录'}
        </p>
      </Card>

      {/* Current Recording Info */}
      {recording && recording.actions.length > 0 && (
        <Card title={`录制数据 (${recording.actions.length} 个动作 / ${(recording.duration ?? 0).toFixed(1)}s)`}>
          <div className="space-y-3">
            <div className="flex gap-2">
              <button
                onClick={onClearRecording}
                className="flex-1 py-2 bg-red-500/10 text-red-400 text-sm rounded-lg border border-red-500/30 hover:bg-red-500/20 transition-colors"
              >
                清除
              </button>
            </div>

            {/* Save */}
            <div className="flex gap-2">
              <input
                type="text"
                value={saveName}
                onChange={e => setSaveName(e.target.value)}
                onKeyDown={e => e.key === 'Enter' && onSaveRecording()}
                placeholder="输入保存名称..."
                className="flex-1 bg-[#334155] text-gray-200 text-sm rounded-lg px-3 py-2 border border-[#475569] focus:border-indigo-500 focus:outline-none"
              />
              <button
                onClick={onSaveRecording}
                disabled={!saveName.trim()}
                className="px-4 py-2 bg-indigo-500 hover:bg-indigo-600 disabled:bg-[#475569] disabled:text-gray-500 text-white text-sm rounded-lg transition-colors"
              >
                保存
              </button>
            </div>

            {/* Action List */}
            <div className="bg-[#0f172a] rounded-lg p-2 max-h-48 overflow-y-auto font-mono text-xs">
              {recording.actions.map((a, i) => (
                <div key={i} className="py-0.5 flex gap-2 items-center">
                  <span className="text-gray-500 w-14 shrink-0">{a.time_offset.toFixed(1)}s</span>
                  <span className={`shrink-0 w-10 ${typeColor(a.type)}`}>{actionTypeLabel(a.type)}</span>
                  <span className="text-gray-300 truncate">{actionDetail(a)}</span>
                </div>
              ))}
            </div>
          </div>
        </Card>
      )}

      {/* Saved Files */}
      <Card title="已保存的录制">
        {savedFiles.length === 0 ? (
          <p className="text-sm text-gray-500 text-center py-4">暂无已保存的录制文件</p>
        ) : (
          <div className="space-y-2">
            {savedFiles.map(name => (
              <div key={name} className="flex items-center justify-between bg-[#0f172a] rounded-lg px-3 py-2">
                <span className="text-sm text-gray-300 font-mono">{name}</span>
                <div className="flex gap-2">
                  <button
                    onClick={() => onLoadRecording(name)}
                    className="px-3 py-1 bg-indigo-500/20 text-indigo-400 text-xs rounded hover:bg-indigo-500/30 transition-colors"
                  >
                    加载
                  </button>
                  <button
                    onClick={() => onDeleteRecording(name)}
                    className="px-3 py-1 bg-red-500/10 text-red-400 text-xs rounded hover:bg-red-500/20 transition-colors"
                  >
                    删除
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </Card>
    </>
  );
}

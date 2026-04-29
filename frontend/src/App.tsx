import { useState, useEffect, useRef } from 'react';
import { GetConfig, UpdateConfig, GetStatus, Start, Stop, GetLogs } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime/runtime';

interface Config {
  move_enabled: boolean;
  move_interval: number;
  move_random: boolean;
  click_enabled: boolean;
  click_interval: number;
  click_type: string;
  click_count: number;
  scroll_enabled: boolean;
  scroll_interval: number;
  scroll_dir: string;
  scroll_amount: number;
  auto_start: boolean;
  minimize_to_tray: boolean;
}

interface LogEntry {
  time: string;
  level: string;
  message: string;
}

type Status = 'stopped' | 'running';

const defaultConfig: Config = {
  move_enabled: false,
  move_interval: 5,
  move_random: true,
  click_enabled: false,
  click_interval: 3,
  click_type: 'left',
  click_count: 1,
  scroll_enabled: false,
  scroll_interval: 5,
  scroll_dir: 'down',
  scroll_amount: 3,
  auto_start: false,
  minimize_to_tray: true,
};

export default function App() {
  const [config, setConfig] = useState<Config>(defaultConfig);
  const [status, setStatus] = useState<Status>('stopped');
  const [activeTab, setActiveTab] = useState<'move' | 'click' | 'scroll' | 'system' | 'logs'>('move');
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const logsEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    loadConfig();
    loadStatus();
    loadLogs();

    const cleanupStatus = EventsOn('statusChanged', (s: string) => {
      setStatus(s as Status);
    });

    const cleanupLog = EventsOn('newLog', (entry: LogEntry) => {
      setLogs(prev => [...prev, entry]);
    });

    return () => {
      cleanupStatus();
      cleanupLog();
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

  async function updateConfig(key: keyof Config, value: any) {
    const newConfig = { ...config, [key]: value };
    setConfig(newConfig);
    try {
      await UpdateConfig(newConfig);
    } catch (e) {
      console.error('Failed to save config:', e);
    }
  }

  const statusColors = {
    stopped: 'bg-gray-500',
    running: 'bg-green-500',
  };

  const statusLabels = {
    stopped: '已停止',
    running: '运行中',
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
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-hidden flex flex-col">
        {/* Tabs */}
        <div className="flex border-b border-[#334155]">
          {[
            { key: 'move', label: '鼠标移动' },
            { key: 'click', label: '鼠标点击' },
            { key: 'scroll', label: '滚轮滚动' },
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
          {activeTab === 'move' && (
            <MoveSettings config={config} updateConfig={updateConfig} />
          )}
          {activeTab === 'click' && (
            <ClickSettings config={config} updateConfig={updateConfig} />
          )}
          {activeTab === 'scroll' && (
            <ScrollSettings config={config} updateConfig={updateConfig} />
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
              <kbd className="px-2 py-1 bg-[#334155] rounded text-indigo-400 font-mono text-xs">F6</kbd>
              <span className="text-gray-400">开始</span>
            </div>
            <div className="flex items-center gap-2">
              <kbd className="px-2 py-1 bg-[#334155] rounded text-indigo-400 font-mono text-xs">F7</kbd>
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
  return (
    <div className="space-y-2">
      <div className="flex justify-between">
        <span className="text-sm text-gray-300">{label}</span>
        <span className="text-sm text-indigo-400 font-medium">{value}{unit}</span>
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

function MoveSettings({ config, updateConfig }: { config: Config; updateConfig: (k: keyof Config, v: any) => void }) {
  return (
    <>
      <Card title="基本设置">
        <Toggle
          label="启用鼠标移动"
          checked={config.move_enabled}
          onChange={v => updateConfig('move_enabled', v)}
        />
      </Card>
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
    </>
  );
}

function ClickSettings({ config, updateConfig }: { config: Config; updateConfig: (k: keyof Config, v: any) => void }) {
  return (
    <>
      <Card title="基本设置">
        <Toggle
          label="启用鼠标点击"
          checked={config.click_enabled}
          onChange={v => updateConfig('click_enabled', v)}
        />
      </Card>
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
    </>
  );
}

function ScrollSettings({ config, updateConfig }: { config: Config; updateConfig: (k: keyof Config, v: any) => void }) {
  return (
    <>
      <Card title="基本设置">
        <Toggle
          label="启用滚轮滚动"
          checked={config.scroll_enabled}
          onChange={v => updateConfig('scroll_enabled', v)}
        />
      </Card>
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
    </>
  );
}

function SystemSettings({ config, updateConfig }: { config: Config; updateConfig: (k: keyof Config, v: any) => void }) {
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
      <Card title="关于">
        <div className="text-center space-y-2 py-2">
          <div className="w-12 h-12 bg-indigo-500 rounded-xl mx-auto flex items-center justify-center">
            <svg className="w-8 h-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
            </svg>
          </div>
          <p className="text-sm font-medium text-gray-200">MousePaw v1.0.0</p>
          <p className="text-xs text-gray-400">鼠标自动化工具</p>
          <p className="text-xs text-gray-500">按 F6 开始，按 F7 停止</p>
        </div>
      </Card>
    </>
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

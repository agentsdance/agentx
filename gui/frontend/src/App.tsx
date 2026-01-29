import { useState, useEffect, useRef } from 'react';
import { Server, Wrench, Package, Terminal, MoreHorizontal, Sun, Moon, Monitor, Download } from 'lucide-react';
import './App.css';
import {
  GetAgents,
  GetMCPMatrix,
  GetSkillsMatrix,
  GetPluginsMatrix,
  InstallMCP,
  RemoveMCP,
  InstallSkillForAgent,
  RemoveSkillFromAgent,
  InstallPluginForAgent,
  RemovePluginFromAgent,
  GetVersion,
  InstallMCPForAll,
  InstallSkillForAll,
  InstallPluginForAll,
} from '../wailsjs/go/main/App';

type AgentInfo = {
  name: string;
  configPath: string;
  exists: boolean;
};

type MCPStatus = {
  name: string;
  description: string;
  agents: Record<string, string>;
};

type SkillStatus = {
  name: string;
  description: string;
  source: string;
  agents: Record<string, string>;
};

type PluginStatus = {
  name: string;
  description: string;
  source: string;
  agents: Record<string, string>;
};

type Tab = 'agents' | 'mcp' | 'skills' | 'plugins';
type Theme = 'light' | 'dark' | 'system';

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('mcp');
  const [agents, setAgents] = useState<AgentInfo[]>([]);
  const [mcpMatrix, setMcpMatrix] = useState<MCPStatus[]>([]);
  const [skillsMatrix, setSkillsMatrix] = useState<SkillStatus[]>([]);
  const [pluginsMatrix, setPluginsMatrix] = useState<PluginStatus[]>([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [version, setVersion] = useState('');
  const [theme, setTheme] = useState<Theme>(() => {
    const saved = localStorage.getItem('agentx-theme');
    return (saved as Theme) || 'light';
  });
  const [showSettings, setShowSettings] = useState(false);
  const settingsRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    loadAgents();
    loadVersion();
  }, []);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (settingsRef.current && !settingsRef.current.contains(event.target as Node)) {
        setShowSettings(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  useEffect(() => {
    const root = document.documentElement;
    localStorage.setItem('agentx-theme', theme);
    
    if (theme === 'system') {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      root.setAttribute('data-theme', prefersDark ? 'dark' : 'light');
    } else {
      root.setAttribute('data-theme', theme);
    }
  }, [theme]);

  useEffect(() => {
    if (activeTab === 'mcp') {
      loadMCPMatrix();
    } else if (activeTab === 'skills') {
      loadSkillsMatrix();
    } else if (activeTab === 'plugins') {
      loadPluginsMatrix();
    }
  }, [activeTab]);

  const loadAgents = async () => {
    const data: AgentInfo[] = await GetAgents();
    setAgents(data);
  };

  const loadVersion = async () => {
    const v = await GetVersion();
    setVersion(v);
  };

  const loadMCPMatrix = async () => {
    setLoading(true);
    try {
      const data = await GetMCPMatrix();
      setMcpMatrix(data || []);
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  };

  const loadSkillsMatrix = async () => {
    setLoading(true);
    try {
      const data = await GetSkillsMatrix();
      setSkillsMatrix(data || []);
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  };

  const loadPluginsMatrix = async () => {
    setLoading(true);
    try {
      const data = await GetPluginsMatrix();
      setPluginsMatrix(data || []);
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  };

  const handleMCPAction = async (agentName: string, mcpName: string, status: string) => {
    setLoading(true);
    try {
      if (status === 'installed') {
        await RemoveMCP(agentName, mcpName);
        setMessage(`Removed ${mcpName} from ${agentName}`);
      } else if (status === 'not_installed') {
        await InstallMCP(agentName, mcpName);
        setMessage(`Installed ${mcpName} to ${agentName}`);
      }
      await loadMCPMatrix();
    } catch (e) {
      console.error(e);
      setMessage(`Error: ${e}`);
    }
    setLoading(false);
    setTimeout(() => setMessage(''), 3000);
  };

  const handleSkillAction = async (agentName: string, skillName: string, source: string, status: string) => {
    setLoading(true);
    try {
      if (status === 'installed') {
        await RemoveSkillFromAgent(agentName, skillName);
        setMessage(`Removed ${skillName} from ${agentName}`);
      } else if (status === 'not_installed') {
        await InstallSkillForAgent(agentName, skillName, source);
        setMessage(`Installed ${skillName} to ${agentName}`);
      }
      await loadSkillsMatrix();
    } catch (e) {
      console.error(e);
      setMessage(`Error: ${e}`);
    }
    setLoading(false);
    setTimeout(() => setMessage(''), 3000);
  };

  const handlePluginAction = async (agentName: string, pluginName: string, source: string, status: string) => {
    setLoading(true);
    try {
      if (status === 'installed') {
        await RemovePluginFromAgent(agentName, pluginName);
        setMessage(`Removed ${pluginName} from ${agentName}`);
      } else if (status === 'not_installed') {
        await InstallPluginForAgent(agentName, pluginName, source);
        setMessage(`Installed ${pluginName} to ${agentName}`);
      }
      await loadPluginsMatrix();
    } catch (e) {
      console.error(e);
      setMessage(`Error: ${e}`);
    }
    setLoading(false);
    setTimeout(() => setMessage(''), 3000);
  };

  const handleInstallMCPAll = async (mcpName: string) => {
    setLoading(true);
    try {
      await InstallMCPForAll(mcpName);
      setMessage(`Installed ${mcpName} to all agents`);
      await loadMCPMatrix();
    } catch (e) {
      console.error(e);
      setMessage(`Error: ${e}`);
    }
    setLoading(false);
    setTimeout(() => setMessage(''), 3000);
  };

  const handleInstallSkillAll = async (skillName: string, source: string) => {
    setLoading(true);
    try {
      await InstallSkillForAll(skillName, source);
      setMessage(`Installed ${skillName} to all agents`);
      await loadSkillsMatrix();
    } catch (e) {
      console.error(e);
      setMessage(`Error: ${e}`);
    }
    setLoading(false);
    setTimeout(() => setMessage(''), 3000);
  };

  const handleInstallPluginAll = async (pluginName: string, source: string) => {
    setLoading(true);
    try {
      await InstallPluginForAll(pluginName, source);
      setMessage(`Installed ${pluginName} to all agents`);
      await loadPluginsMatrix();
    } catch (e) {
      console.error(e);
      setMessage(`Error: ${e}`);
    }
    setLoading(false);
    setTimeout(() => setMessage(''), 3000);
  };

  const getStatusClass = (status: string) => {
    switch (status) {
      case 'installed':
        return 'status-installed';
      case 'not_installed':
        return 'status-not-installed';
      case 'n/a':
        return 'status-na';
      case 'error':
        return 'status-error';
      default:
        return '';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'installed':
        return '✓ installed';
      case 'not_installed':
        return '○ ---';
      case 'n/a':
        return '○ n/a';
      case 'error':
        return '✗ error';
      default:
        return status;
    }
  };

  const formatMCPName = (name: string) => {
    switch (name) {
      case 'playwright':
        return 'Playwright';
      case 'context7':
        return 'Context7';
      case 'remix-icon':
        return 'Remix Icon';
      default:
        return name;
    }
  };

  return (
    <div className="app">
      <header className="header">
        <div className="header-content">
          <h1>AgentX</h1>
          {version && <span className="version-tag">{version}</span>}
        </div>
        <p className="subtitle">Manage MCP servers, skills, and plugins for AI coding agents</p>
      </header>

      <nav className="tabs">
        <button
          className={`tab ${activeTab === 'mcp' ? 'active' : ''}`}
          onClick={() => setActiveTab('mcp')}
        >
          <Server size={16} /> MCP Servers
        </button>
        <button
          className={`tab ${activeTab === 'skills' ? 'active' : ''}`}
          onClick={() => setActiveTab('skills')}
        >
          <Wrench size={16} /> Skills
        </button>
        <button
          className={`tab ${activeTab === 'plugins' ? 'active' : ''}`}
          onClick={() => setActiveTab('plugins')}
        >
          <Package size={16} /> Plugins
        </button>
        <button
          className={`tab ${activeTab === 'agents' ? 'active' : ''}`}
          onClick={() => setActiveTab('agents')}
        >
          <Terminal size={16} /> Code Agents
        </button>
      </nav>

      {message && <div className="message">{message}</div>}

      <main className="content">
        {activeTab === 'agents' && (
          <div className="panel">
            <h2>AI Coding Agents</h2>
            <div className="agents-grid">
              {agents.map((agent) => (
                <div
                  key={agent.name}
                  className={`agent-card ${agent.exists ? 'installed' : 'not-installed'}`}
                >
                  <div className="agent-header">
                    <h3>{agent.name}</h3>
                    <span className={`status ${agent.exists ? 'active' : 'inactive'}`}>
                      {agent.exists ? 'Installed' : 'Not Found'}
                    </span>
                  </div>
                  <p className="config-path">{agent.configPath}</p>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === 'mcp' && (
          <div className="panel">
            <h2>MCP Server Status</h2>
            {loading ? (
              <div className="loading">Loading...</div>
            ) : (
              <div className="matrix-container">
                <table className="matrix-table">
                  <thead>
                    <tr>
                      <th className="row-header">MCP Server</th>
                      {agents.map((agent) => (
                        <th key={agent.name} className="col-header">
                          {agent.name}
                        </th>
                      ))}
                      <th className="col-header action-header">Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    {mcpMatrix.map((mcp) => (
                      <tr key={mcp.name}>
                        <td className="row-label">
                          <span className="item-name">{formatMCPName(mcp.name)}</span>
                          <span className="item-desc">{mcp.description}</span>
                        </td>
                        {agents.map((agent) => {
                          const status = mcp.agents[agent.name] || 'n/a';
                          const isClickable = status === 'installed' || status === 'not_installed';
                          return (
                            <td
                              key={agent.name}
                              className={`matrix-cell ${getStatusClass(status)} ${isClickable ? 'clickable' : ''}`}
                              onClick={() => isClickable && handleMCPAction(agent.name, mcp.name, status)}
                              title={isClickable ? (status === 'installed' ? 'Click to remove' : 'Click to install') : ''}
                            >
                              {getStatusText(status)}
                            </td>
                          );
                        })}
                        <td className="matrix-cell action-cell">
                          <button
                            className="install-all-btn"
                            onClick={() => handleInstallMCPAll(mcp.name)}
                            title={`Install ${mcp.name} to all agents`}
                          >
                            <Download size={14} /> Install All
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
            <p className="hint">Click on a cell to install/remove</p>
          </div>
        )}

        {activeTab === 'skills' && (
          <div className="panel">
            <h2>Skills Status</h2>
            {loading ? (
              <div className="loading">Loading...</div>
            ) : skillsMatrix.length === 0 ? (
              <p className="empty">No skills available in registry</p>
            ) : (
              <div className="matrix-container">
                <table className="matrix-table">
                  <thead>
                    <tr>
                      <th className="row-header">Skill</th>
                      {agents.map((agent) => (
                        <th key={agent.name} className="col-header">
                          {agent.name}
                        </th>
                      ))}
                      <th className="col-header action-header">Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    {skillsMatrix.map((skill) => (
                      <tr key={skill.name}>
                        <td className="row-label">
                          <span className="item-name">{skill.name}</span>
                          <span className="item-desc">{skill.description}</span>
                        </td>
                        {agents.map((agent) => {
                          const status = skill.agents[agent.name] || 'n/a';
                          const isClickable = status === 'installed' || status === 'not_installed';
                          return (
                            <td
                              key={agent.name}
                              className={`matrix-cell ${getStatusClass(status)} ${isClickable ? 'clickable' : ''}`}
                              onClick={() => isClickable && handleSkillAction(agent.name, skill.name, skill.source, status)}
                              title={isClickable ? (status === 'installed' ? 'Click to remove' : 'Click to install') : ''}
                            >
                              {getStatusText(status)}
                            </td>
                          );
                        })}
                        <td className="matrix-cell action-cell">
                          <button
                            className="install-all-btn"
                            onClick={() => handleInstallSkillAll(skill.name, skill.source)}
                            title={`Install ${skill.name} to all agents`}
                          >
                            <Download size={14} /> Install All
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
            <p className="hint">Click on a cell to install/remove</p>
          </div>
        )}

        {activeTab === 'plugins' && (
          <div className="panel">
            <h2>Plugins Status</h2>
            {loading ? (
              <div className="loading">Loading...</div>
            ) : pluginsMatrix.length === 0 ? (
              <p className="empty">No plugins available in registry</p>
            ) : (
              <div className="matrix-container">
                <table className="matrix-table">
                  <thead>
                    <tr>
                      <th className="row-header">Plugin</th>
                      {agents.map((agent) => (
                        <th key={agent.name} className="col-header">
                          {agent.name}
                        </th>
                      ))}
                      <th className="col-header action-header">Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    {pluginsMatrix.map((plugin) => (
                      <tr key={plugin.name}>
                        <td className="row-label">
                          <span className="item-name">{plugin.name}</span>
                          <span className="item-desc">{plugin.description}</span>
                        </td>
                        {agents.map((agent) => {
                          const status = plugin.agents[agent.name] || 'n/a';
                          const isClickable = status === 'installed' || status === 'not_installed';
                          return (
                            <td
                              key={agent.name}
                              className={`matrix-cell ${getStatusClass(status)} ${isClickable ? 'clickable' : ''}`}
                              onClick={() => isClickable && handlePluginAction(agent.name, plugin.name, plugin.source, status)}
                              title={isClickable ? (status === 'installed' ? 'Click to remove' : 'Click to install') : ''}
                            >
                              {getStatusText(status)}
                            </td>
                          );
                        })}
                        <td className="matrix-cell action-cell">
                          <button
                            className="install-all-btn"
                            onClick={() => handleInstallPluginAll(plugin.name, plugin.source)}
                            title={`Install ${plugin.name} to all agents`}
                          >
                            <Download size={14} /> Install All
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
            <p className="hint">Click on a cell to install/remove</p>
          </div>
        )}
      </main>

      <div className="settings-corner" ref={settingsRef}>
        <button
          className={`settings-trigger ${showSettings ? 'active' : ''}`}
          onClick={() => setShowSettings(!showSettings)}
          title="Settings"
        >
          <MoreHorizontal size={20} />
        </button>
        {showSettings && (
          <div className="settings-popup">
            <div className="settings-popup-header">Theme</div>
            <button
              className={`settings-popup-item ${theme === 'light' ? 'active' : ''}`}
              onClick={() => setTheme('light')}
            >
              <Sun size={16} />
              <span>Light</span>
            </button>
            <button
              className={`settings-popup-item ${theme === 'dark' ? 'active' : ''}`}
              onClick={() => setTheme('dark')}
            >
              <Moon size={16} />
              <span>Dark</span>
            </button>
            <button
              className={`settings-popup-item ${theme === 'system' ? 'active' : ''}`}
              onClick={() => setTheme('system')}
            >
              <Monitor size={16} />
              <span>System</span>
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;

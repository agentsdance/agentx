export namespace main {
	
	export class AgentInfo {
	    name: string;
	    configPath: string;
	    exists: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AgentInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.configPath = source["configPath"];
	        this.exists = source["exists"];
	    }
	}
	export class MCPInfo {
	    name: string;
	    config: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new MCPInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.config = source["config"];
	    }
	}
	export class MCPStatus {
	    name: string;
	    description: string;
	    agents: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new MCPStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.agents = source["agents"];
	    }
	}
	export class PluginInfo {
	    name: string;
	    description: string;
	    source: string;
	    version: string;
	    author: string;
	    components: string[];
	
	    static createFrom(source: any = {}) {
	        return new PluginInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.source = source["source"];
	        this.version = source["version"];
	        this.author = source["author"];
	        this.components = source["components"];
	    }
	}
	export class PluginStatus {
	    name: string;
	    description: string;
	    source: string;
	    agents: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new PluginStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.source = source["source"];
	        this.agents = source["agents"];
	    }
	}
	export class SkillInfo {
	    name: string;
	    description: string;
	    author: string;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new SkillInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.author = source["author"];
	        this.source = source["source"];
	    }
	}
	export class SkillStatus {
	    name: string;
	    description: string;
	    source: string;
	    agents: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new SkillStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.source = source["source"];
	        this.agents = source["agents"];
	    }
	}

}


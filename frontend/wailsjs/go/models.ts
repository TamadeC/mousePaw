export namespace config {
	
	export class HotkeyConfig {
	    start: string;
	    stop: string;
	    pause: string;
	
	    static createFrom(source: any = {}) {
	        return new HotkeyConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.start = source["start"];
	        this.stop = source["stop"];
	        this.pause = source["pause"];
	    }
	}
	export class Config {
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
	    hotkeys: HotkeyConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.operation_type = source["operation_type"];
	        this.move_interval = source["move_interval"];
	        this.move_random = source["move_random"];
	        this.click_interval = source["click_interval"];
	        this.click_type = source["click_type"];
	        this.click_count = source["click_count"];
	        this.scroll_interval = source["scroll_interval"];
	        this.scroll_dir = source["scroll_dir"];
	        this.scroll_amount = source["scroll_amount"];
	        this.type_interval = source["type_interval"];
	        this.type_text = source["type_text"];
	        this.replay_interval = source["replay_interval"];
	        this.replay_repeat = source["replay_repeat"];
	        this.replay_file = source["replay_file"];
	        this.auto_start = source["auto_start"];
	        this.minimize_to_tray = source["minimize_to_tray"];
	        this.hotkeys = this.convertValues(source["hotkeys"], HotkeyConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace log {
	
	export class LogEntry {
	    time: string;
	    level: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new LogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.level = source["level"];
	        this.message = source["message"];
	    }
	}

}

export namespace recorder {
	
	export class RecordedAction {
	    time_offset: number;
	    type: string;
	    x?: number;
	    y?: number;
	    button?: string;
	    direction?: string;
	    amount?: number;
	    keychar?: string;
	
	    static createFrom(source: any = {}) {
	        return new RecordedAction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time_offset = source["time_offset"];
	        this.type = source["type"];
	        this.x = source["x"];
	        this.y = source["y"];
	        this.button = source["button"];
	        this.direction = source["direction"];
	        this.amount = source["amount"];
	        this.keychar = source["keychar"];
	    }
	}
	export class RecorderStatusInfo {
	    status: string;
	    count: number;
	    duration: number;
	
	    static createFrom(source: any = {}) {
	        return new RecorderStatusInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.count = source["count"];
	        this.duration = source["duration"];
	    }
	}
	export class Recording {
	    actions: RecordedAction[];
	    duration: number;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Recording(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.actions = this.convertValues(source["actions"], RecordedAction);
	        this.duration = source["duration"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}


export namespace config {
	
	export class Config {
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
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.move_enabled = source["move_enabled"];
	        this.move_interval = source["move_interval"];
	        this.move_random = source["move_random"];
	        this.click_enabled = source["click_enabled"];
	        this.click_interval = source["click_interval"];
	        this.click_type = source["click_type"];
	        this.click_count = source["click_count"];
	        this.scroll_enabled = source["scroll_enabled"];
	        this.scroll_interval = source["scroll_interval"];
	        this.scroll_dir = source["scroll_dir"];
	        this.scroll_amount = source["scroll_amount"];
	        this.auto_start = source["auto_start"];
	        this.minimize_to_tray = source["minimize_to_tray"];
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


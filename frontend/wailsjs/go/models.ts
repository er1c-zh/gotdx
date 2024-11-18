export namespace api {
	
	export enum MsgKey {
	    init = "init",
	    processMsg = "processMsg",
	    connectionStatus = "connectionStatus",
	}

}

export namespace models {
	
	export enum MarketType {
	    深圳 = 0x0,
	    上海 = 0x1,
	    北京 = 0x2,
	}
	export class ProcessInfo {
	    Type: number;
	    Msg: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Type = source["Type"];
	        this.Msg = source["Msg"];
	    }
	}

}


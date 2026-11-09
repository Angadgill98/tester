export namespace models {
	
	export class ActionInfo {
	    Id: number;
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new ActionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	    }
	}
	export class ChildApplicationInfo {
	    Id: number;
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new ChildApplicationInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	    }
	}
	export class SequenceInfo {
	    Id: number;
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new SequenceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	    }
	}
	export class ApplicationFrontend {
	    AppId: number;
	    AppName: string;
	    ParentAppId: number;
	    EnvironmentId: sql.NullInt64;
	    Actions: ActionInfo[];
	    Sequences: SequenceInfo[];
	    Children: ChildApplicationInfo[];
	
	    static createFrom(source: any = {}) {
	        return new ApplicationFrontend(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AppId = source["AppId"];
	        this.AppName = source["AppName"];
	        this.ParentAppId = source["ParentAppId"];
	        this.EnvironmentId = this.convertValues(source["EnvironmentId"], sql.NullInt64);
	        this.Actions = this.convertValues(source["Actions"], ActionInfo);
	        this.Sequences = this.convertValues(source["Sequences"], SequenceInfo);
	        this.Children = this.convertValues(source["Children"], ChildApplicationInfo);
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

export namespace sql {
	
	export class NullInt64 {
	    Int64: number;
	    Valid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NullInt64(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Int64 = source["Int64"];
	        this.Valid = source["Valid"];
	    }
	}

}


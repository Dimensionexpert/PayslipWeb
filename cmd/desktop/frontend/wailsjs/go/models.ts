export namespace models {
	
	export class Cluster {
	    ID: number;
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new Cluster(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	    }
	}
	export class Employee {
	    ID: number;
	    SchoolID: number;
	    SerialNo: string;
	    Name: string;
	    ShalarthID: string;
	    Gender: string;
	    Designation: string;
	    GPFNo: string;
	    DCPSNo: string;
	    PRANNo: string;
	    PAN: string;
	    Aadhaar: string;
	    Mobile: string;
	    Email: string;
	    DDOBankName: string;
	    DDOBankAccount: string;
	    DDOBankIFSC: string;
	    BankName: string;
	    BankAccount: string;
	    BankIFSC: string;
	    BranchName: string;
	    PayMatrix: string;
	    UDISECode: string;
	
	    static createFrom(source: any = {}) {
	        return new Employee(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.SchoolID = source["SchoolID"];
	        this.SerialNo = source["SerialNo"];
	        this.Name = source["Name"];
	        this.ShalarthID = source["ShalarthID"];
	        this.Gender = source["Gender"];
	        this.Designation = source["Designation"];
	        this.GPFNo = source["GPFNo"];
	        this.DCPSNo = source["DCPSNo"];
	        this.PRANNo = source["PRANNo"];
	        this.PAN = source["PAN"];
	        this.Aadhaar = source["Aadhaar"];
	        this.Mobile = source["Mobile"];
	        this.Email = source["Email"];
	        this.DDOBankName = source["DDOBankName"];
	        this.DDOBankAccount = source["DDOBankAccount"];
	        this.DDOBankIFSC = source["DDOBankIFSC"];
	        this.BankName = source["BankName"];
	        this.BankAccount = source["BankAccount"];
	        this.BankIFSC = source["BankIFSC"];
	        this.BranchName = source["BranchName"];
	        this.PayMatrix = source["PayMatrix"];
	        this.UDISECode = source["UDISECode"];
	    }
	}
	export class PayslipRecord {
	    Month: number;
	    Year: number;
	    ID: number;
	    EmployeeID: number;
	    BasicPay: number;
	    DA: number;
	    HRA: number;
	    HRAArrear: number;
	    TA: number;
	    TAArrear: number;
	    TribalAllowance: number;
	    WashingAllowance: number;
	    DAArrears: number;
	    BasicArrears: number;
	    CLA: number;
	    NPSEmprAllow: number;
	    TotalPay: number;
	    FA: number;
	    GrossAfterFA: number;
	    GPF: number;
	    GPFAdvance: number;
	    PT: number;
	    GISZP: number;
	    GISScout: number;
	    DCPSRegular: number;
	    DCPSDelayed: number;
	    DCPSPayArrears: number;
	    RevenueStamp: number;
	    DCPSDAArrears: number;
	    GroupAccidentalPolicy: number;
	    NAA: number;
	    TotalGovtDeductions: number;
	    GrossAfterGovtDeductions: number;
	    NPSEmprContri: number;
	    NPSEmpContri: number;
	    NPSEmprContriArr: number;
	    NPSEmpContriArr: number;
	    NPSTotal: number;
	    GrossAfterNPS: number;
	    IncomeTax: number;
	    CoopBank: number;
	    NGRLIC: number;
	    NGRSocietyLoan: number;
	    NGRMisc: number;
	    NGROtherRecovery: number;
	    NGRRD: number;
	    NGROtherDeduction: number;
	    NGRTotalDeduction: number;
	    EmployeeNetSalary: number;
	    Remarks: string;
	    UDISECode: string;
	    ShalarthID: string;
	
	    static createFrom(source: any = {}) {
	        return new PayslipRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Month = source["Month"];
	        this.Year = source["Year"];
	        this.ID = source["ID"];
	        this.EmployeeID = source["EmployeeID"];
	        this.BasicPay = source["BasicPay"];
	        this.DA = source["DA"];
	        this.HRA = source["HRA"];
	        this.HRAArrear = source["HRAArrear"];
	        this.TA = source["TA"];
	        this.TAArrear = source["TAArrear"];
	        this.TribalAllowance = source["TribalAllowance"];
	        this.WashingAllowance = source["WashingAllowance"];
	        this.DAArrears = source["DAArrears"];
	        this.BasicArrears = source["BasicArrears"];
	        this.CLA = source["CLA"];
	        this.NPSEmprAllow = source["NPSEmprAllow"];
	        this.TotalPay = source["TotalPay"];
	        this.FA = source["FA"];
	        this.GrossAfterFA = source["GrossAfterFA"];
	        this.GPF = source["GPF"];
	        this.GPFAdvance = source["GPFAdvance"];
	        this.PT = source["PT"];
	        this.GISZP = source["GISZP"];
	        this.GISScout = source["GISScout"];
	        this.DCPSRegular = source["DCPSRegular"];
	        this.DCPSDelayed = source["DCPSDelayed"];
	        this.DCPSPayArrears = source["DCPSPayArrears"];
	        this.RevenueStamp = source["RevenueStamp"];
	        this.DCPSDAArrears = source["DCPSDAArrears"];
	        this.GroupAccidentalPolicy = source["GroupAccidentalPolicy"];
	        this.NAA = source["NAA"];
	        this.TotalGovtDeductions = source["TotalGovtDeductions"];
	        this.GrossAfterGovtDeductions = source["GrossAfterGovtDeductions"];
	        this.NPSEmprContri = source["NPSEmprContri"];
	        this.NPSEmpContri = source["NPSEmpContri"];
	        this.NPSEmprContriArr = source["NPSEmprContriArr"];
	        this.NPSEmpContriArr = source["NPSEmpContriArr"];
	        this.NPSTotal = source["NPSTotal"];
	        this.GrossAfterNPS = source["GrossAfterNPS"];
	        this.IncomeTax = source["IncomeTax"];
	        this.CoopBank = source["CoopBank"];
	        this.NGRLIC = source["NGRLIC"];
	        this.NGRSocietyLoan = source["NGRSocietyLoan"];
	        this.NGRMisc = source["NGRMisc"];
	        this.NGROtherRecovery = source["NGROtherRecovery"];
	        this.NGRRD = source["NGRRD"];
	        this.NGROtherDeduction = source["NGROtherDeduction"];
	        this.NGRTotalDeduction = source["NGRTotalDeduction"];
	        this.EmployeeNetSalary = source["EmployeeNetSalary"];
	        this.Remarks = source["Remarks"];
	        this.UDISECode = source["UDISECode"];
	        this.ShalarthID = source["ShalarthID"];
	    }
	}
	export class School {
	    ID: number;
	    UDISECode: string;
	    Name: string;
	    DDOCode: string;
	    Block: string;
	    Cluster: string;
	
	    static createFrom(source: any = {}) {
	        return new School(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.UDISECode = source["UDISECode"];
	        this.Name = source["Name"];
	        this.DDOCode = source["DDOCode"];
	        this.Block = source["Block"];
	        this.Cluster = source["Cluster"];
	    }
	}
	export class PayslipExport {
	    Cluster: Cluster;
	    School: School;
	    Employee: Employee;
	    Payslip: PayslipRecord;
	
	    static createFrom(source: any = {}) {
	        return new PayslipExport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Cluster = this.convertValues(source["Cluster"], Cluster);
	        this.School = this.convertValues(source["School"], School);
	        this.Employee = this.convertValues(source["Employee"], Employee);
	        this.Payslip = this.convertValues(source["Payslip"], PayslipRecord);
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


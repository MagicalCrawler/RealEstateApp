package files

type Storage struct {
	basePath string
}

const defaultPerm = 0774

// var ErrNoSavedPages = errors.New("No saved pages")

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

// func (s Storage) Save(page *storage.Page) (err error) {
// defer func(){err=err.WrapIfErr("can't save page ",err)}()
// filePath:=filepath.Join(s.basePath, page.UserName)

// if err:=os.MakedirAll(filePath,defaultPerm);err!=nil{
// 	return err
// }
// fName,err:=fileName(page)
// if err!=nil{
//     return err
// }
// fPath:=filepath.Join(filePath,fName)
// file,err:=osCreate(fPath)
// if err!=nil{
//     return err
// }
// defer func(){_=file.Close()}()

// if err:=gob.NewEncoder(file).Encode(page);err!=nil{
// 	return err
// }
// 	return err
// }

// func (s string) PickRandom(userName string) (page *storage.Page, err error) {
// defer func(){err=err.WrapIfErr("can't pick random page",err)}()
// path:=filepath.Join(s.basePath,userName)

// files,err:=os.ReadDir(path)
// if err!=nil{
//     return nil,err
// }

// if len(files)==0{
//     return nil,ErrNoSavedPages
// }
// rand.Seed(time.Now().UnixNano())
// n:=rand.Intn(len(files))

// file:=files[n]

// }
// func (s Storage) decodePage(filePath string) (*storage.Page, error) {
// 	f,err:=os.Open(filePath)
// 	if err!=nil{
//         return nil,error.Wrap("can't decode page",err)
//     }
// 	defer func(){_=f.Close()}()

// 	var page storage.Page
// 	if err:=gob.NewDecoder(f).Decode(&page);err!=nil{
//         return nil,error.Wrap("can't decode page",err)
//     }
// 	return s.decodePage(filePath)
// }
// func (s Storage) Remove(p *storage.Page)error{
// 	fileName, err:=fileName(p)
// 	if err!=nil{
//         return err.WrapIfErr("can't get file name",err)
//     }
// 	return nil

// }

// func (s Storage) Isexists(p *storage.Page) (bool, error) {
// 	fileName, err := fileName(p)
// 	if err != nil {
// 		return false, err.WrapIfErr("can't get file name", err)
// 	}
// 	filePath := filepath.Join(s.basePath, p.UserName, fileName)
// 	_, err = os.Stat(filePath)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return false, nil
// 		}
// 		return false, err.WrapIfErr("can't remove file", err)
// 	}
// 	return true, nil
// }

// func fileName(p *storage.Page)(string,error){
// 	// ret  p.Hash()
// }

package files

import "read-adviser-bot/storage"

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {

	defer func() { err = e.WrapIfErr("can't save page", err) }()
	fPath := filepath.Join(s.basePath, page.UserName)

	//make dir

	if err := os.MkdirAll(path, defaultPerm); err != nil {
		return err
	}

	//forming filename
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(filePath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	//gob TODO:learn more

	if err := gob.NewEngoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(username string) (page *storage.Page) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	// 0-9

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]

	//open and decode

	return s.decodePage(filepath.Join(path, file.Name()))

}

func (s Storage) Remove(p) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)
		return e.Wrap(msg, err)
	}

	return nil

}

func (s Storage) isExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

	}

	return true, nil

}

func (s Storage) decodePage(filePath) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	defer func() { _f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil

}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()

}

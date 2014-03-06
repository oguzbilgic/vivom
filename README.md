# Vivom [![Build Status](https://travis-ci.org/oguzbilgic/vivom.png?branch=master)](https://travis-ci.org/oguzbilgic/vivom)

dead simple, tiny, but powerful Go ORM library

## Usage

Implicitly implement the vivom's interfaces as necessary.

```go
type Subscriber struct {
	ID        int
	Name      string
	Email     string
	DateAdded string
}

func (s *Subscriber) GetID() int {
	return s.ID
}

func (s *Subscriber) SetID(ID int) {
	s.ID = ID
}

func (s *Subscriber) Validate() error {
	if s.Name == "" {
		return errors.New("empty name")
	}

	if s.Email == "" {
		return errors.New("invalid email")
	}

	return nil
}

func (s *Subscriber) Table() string {
	return "subscribers"
}

func (s *Subscriber) Columns() []string {
	return []string{"id", "name", "email", "date_added"}
}

func (s *Subscriber) Values() []interface{} {
	return []interface{}{s.Name, s.Email, time.Now().Unix()}
}

func (s *Subscriber) ScanValues() []interface{} {
	return []interface{}{&s.ID, &s.Name, &s.Email, &t.DateAdded}
}
```

Import the vivom package

```go
import "github.com/oguzbilgic/vivom"
```

Now access your records easily using vivom's db functions.

```go
subscriber := &Subscriber{}
err := vivom.Select(subscriber, 23152, db)
if err != nil {
	panic(err)
}

fmt.Println("Subscriber #"+subscriber.ID+" is "+subscriber.Name)
```

You can also insert new records to the database

```go
subscriber := &Subsriber{Name: "John Doe", Email: "foo@bar.com"}
err := vivom.Insert(subscriber, db)
if err != nil {
	panic(err)
}

fmt.Println("ID of the new subscriber is "+subscriber.ID)
```

For managing multiple database records, your struct should also implement vivom.TableRows interface.

## Documentation

http://godoc.org/github.com/oguzbilgic/vivom

## License

The MIT License (MIT)

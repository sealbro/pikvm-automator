package storage

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	NotFound = fmt.Errorf("not found")
)

type Command struct {
	ID          string `yaml:"id"`
	Description string `yaml:"description"`
	Expression  string `yaml:"expression"`
}

type CommandRepository struct {
	commands map[string]Command
	filePath string
}

func NewCommandRepository(filePath string) *CommandRepository {
	return &CommandRepository{
		filePath: filePath,
	}
}

func (r *CommandRepository) GetCommands() ([]Command, error) {
	err := r.load()
	if err != nil {
		return nil, err
	}

	commands := make([]Command, 0, len(r.commands))
	for _, command := range r.commands {
		commands = append(commands, command)
	}

	return commands, nil
}

func (r *CommandRepository) GetCommand(id string) (*Command, error) {
	err := r.load()
	if err != nil {
		return nil, err
	}

	command, ok := r.commands[id]
	if ok {
		return &command, nil
	}

	return nil, NotFound
}

func (r *CommandRepository) CreateCommand(command Command) error {
	err := r.load()
	if err != nil {
		return err
	}

	r.commands[command.ID] = command

	err = r.save()
	if err != nil {
		return err
	}

	return nil
}

func (r *CommandRepository) DeleteCommand(id string) error {
	err := r.load()
	if err != nil {
		return err
	}

	if _, ok := r.commands[id]; ok {
		delete(r.commands, id)
		err = r.save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *CommandRepository) load() error {
	if r.commands != nil {
		return nil
	}

	_, err := os.Stat(r.filePath)
	if os.IsNotExist(err) {
		r.commands = make(map[string]Command)
		return nil
	}

	file, err := os.Open(r.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&r.commands)
	if err != nil {
		return err
	}

	return nil
}

func (r *CommandRepository) save() error {
	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(r.commands)
	if err != nil {
		return err
	}

	return nil
}

package flow_dsl

import (
	"encoding/json"
	"github.com/dapr-platform/common"
	"github.com/pkg/errors"
	"workflow-service/entity"
)

type JsonParser struct {
}

type WorkflowRule struct {
	Cells []map[string]any `json:"cells"`
}

func (p *JsonParser) Parse(content string) (workflow Workflow, err error) {
	var rule WorkflowRule
	err = json.Unmarshal([]byte(content), &rule)
	if err != nil {
		err = errors.Wrap(err, "failed to parse json")
		return
	}
	statementMap := make(map[string]*Statement, 0)
	for _, cellMap := range rule.Cells { //先处理卡片，再处理连线
		shape := cellMap["shape"]
		if shape != "edge" {
			var cell entity.Cell
			err = cell.FromMap(cellMap)
			if err != nil {
				common.Logger.Error("cell error , can't parse json ", err)
				return
			}

			statementMap[cell.Id] = &Statement{
				Id:               cell.Id,
				Type:             cell.Shape,
				Properties:       cell.Data,
				IncomingBusiness: make(map[string]map[string]string, 0),
				OutgoingBusiness: make(map[string]map[string]string, 0),
				Activity: &ActivityInvocation{
					Id:   cell.Id,
					Name: cell.Shape,
				},
				NextStatements: make([]*Statement, 0),
			}
		}

	}

	for _, cellMap := range rule.Cells {
		shape := cellMap["shape"]
		if shape == "edge" {
			var sStatement *Statement
			var sBusiness, tBusiness string
			var tStatement *Statement

			var edge entity.Edge
			err = edge.FromMap(cellMap)
			if err != nil {
				common.Logger.Error("edge error , can't parse json ", err)
				err = errors.Wrap(err, "failed to parse edge json")
				return
			}
			sStatement, ok := statementMap[edge.Source.Cell]
			if !ok {
				err = errors.New("can't find source statement by edge source cell " + edge.Source.Cell)
				return
			}
			sBusiness = edge.Source.Port
			tStatement, ok = statementMap[edge.Target.Cell]
			if !ok {
				err = errors.New("can't find target statement by edge target cell " + edge.Target.Cell)
				return
			}

			if sStatement == nil || tStatement == nil {
				common.Logger.Error("edge error , can't find source or target statement s=", sStatement, " t=", tStatement)
				continue
			}
			tBusiness = edge.Target.Port
			if sStatement.OutgoingBusiness[sBusiness] == nil {
				sStatement.OutgoingBusiness[sBusiness] = make(map[string]string)
			}
			sStatement.OutgoingBusiness[sBusiness][tBusiness] = tBusiness
			if tStatement.IncomingBusiness[tBusiness] == nil {
				tStatement.IncomingBusiness[tBusiness] = make(map[string]string)
			}
			tStatement.IncomingBusiness[tBusiness][sBusiness] = sBusiness
			shouldIgnore := false
			for _, t := range sStatement.NextStatements {
				if t.Id == tStatement.Id {
					common.Logger.Warn("statement already exist,ignore it ", tStatement)
					shouldIgnore = true
					break
				}
			}
			if !shouldIgnore {
				sStatement.NextStatements = append(sStatement.NextStatements, tStatement)
			}

		}
	}

	workflow.Statements = make([]*Statement, 0)
	for _, sta := range statementMap {
		if len(sta.IncomingBusiness) == 0 {
			workflow.Statements = append(workflow.Statements, sta)
		}
	}
	if len(workflow.Statements) == 0 {
		err = errors.New("can't find top nodes,which is no incoming edge  ")
	}

	return
}

package usecases

import (
	"context"

	"github.com/plopyblopy/orgstruct/internal/domain"
)

// PostDepartment usecase для добавления нового Department.
func PostDepartment(repo domain.DepartamentAdder) func(ctx context.Context, name string, parentId *int) (*domain.Department, error) {
	return func(ctx context.Context, name string, parentId *int) (*domain.Department, error) {
		model, err := domain.NewDepartment(name, parentId)
		if err != nil {
			return nil, err
		}

		err = repo.Add(ctx, model)
		if err != nil {
			return nil, err
		}

		return model, nil
	}
}

// PostEmployee usecase для добавления нового Employee.
func PostEmployee(repo domain.EmployeeAdder) func(ctx context.Context, r domain.AddEmployeeRequest) (*domain.Employee, error) {
	return func(ctx context.Context, r domain.AddEmployeeRequest) (*domain.Employee, error) {
		model, err := domain.NewEmployee(r.DepartmentId, r.FullName, r.Position, r.HiredAt)
		if err != nil {
			return nil, err
		}

		err = repo.Add(ctx, model)
		if err != nil {
			return nil, err
		}

		return model, nil
	}
}

// GetDepartmentTree usecase для получения Department с или без Employees и с вложенными дочерними Departments
func GetDepartmentTree(deptRepo domain.WithDepthHierarchyGetter, empRepo domain.EmployeesByDepartmentIdsGetter) func(ctx context.Context, deptId int, depth int, includeEmployees bool) (*domain.DepartmentWithChildResponse, error) {
	return func(ctx context.Context, deptId int, depth int, includeEmployees bool) (*domain.DepartmentWithChildResponse, error) {
		// получение Departments с полем depth.
		flatDepts, err := deptRepo.GetWithDepthHierarchy(ctx, deptId, depth)
		if err != nil {
			return nil, err
		}

		// получение employee на основе departments.
		var empMap map[int][]domain.EmployeeResponse
		if includeEmployees {
			deptIds := make([]int, len(flatDepts))
			for i := range flatDepts {
				deptIds[i] = flatDepts[i].Id
			}

			emps, err := empRepo.GetByDepartmentIds(ctx, deptIds...)
			if err != nil {
				return nil, err
			}

			empMap = make(map[int][]domain.EmployeeResponse, len(flatDepts))
			for i := range emps {
				emp := emps[i]
				empMap[emp.DepartmentId] = append(empMap[emp.DepartmentId], *emp)
			}
		}

		// преобразование в конечный dto.
		treeMap := make(map[int]*domain.DepartmentWithChildResponse, len(flatDepts))
		for i := range flatDepts {
			d := flatDepts[i]
			treeMap[d.Id] = &domain.DepartmentWithChildResponse{
				DepartmentResponse: domain.DepartmentResponse{
					Id:        d.Id,
					Name:      d.Name,
					ParentId:  d.ParentId,
					CreatedAt: d.CreatedAt,
				},
				Employees: empMap[d.Id],
				Children:  []*domain.DepartmentWithChildResponse{},
			}
		}

		// постройка иерархического дерева.
		var root *domain.DepartmentWithChildResponse
		for i := range flatDepts {
			d := flatDepts[i]
			node := treeMap[d.Id]
			if d.ParentId == nil {
				root = node
			} else if parent, ok := treeMap[*d.ParentId]; ok {
				parent.Children = append(parent.Children, node)
			}
		}

		// если parent department не корень
		if root == nil {
			for i := range flatDepts {
				d := flatDepts[i]

				if d.Depth == 0 {
					root = treeMap[d.Id]
					break
				}
			}
		}

		return root, nil
	}
}

// PatchDepartment для частичного обновления Department.
func PatchDepartment(repo domain.DepartmentUpdater) func(ctx context.Context, dept domain.UpdateDepartment) (*domain.Department, error) {
	return func(ctx context.Context, dept domain.UpdateDepartment) (*domain.Department, error) {
		if dept.Name.Valid {
			v := domain.NewValidator()
			v.NotNil("name", dept.Name.Value)
			v.MinMax("name", domain.DeptNameMinLen, domain.DeptNameMaxLen, len(*dept.Name.Value))
			if err := v.Validate(); err != nil {
				return nil, err
			}
		}

		model, err := repo.Update(ctx, dept)
		if err != nil {
			return nil, err
		}

		return model, nil
	}
}
